import { Inject, Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import Bottleneck from 'bottleneck';
import type { HttpClientPort } from '@/core/application/ports/http-client.port';
import type { ProxyManagerPort } from '@/core/application/ports/proxy-manager.port';
import {
  HTTP_CLIENT_PORT,
  PROXY_MANAGER_PORT,
} from '@/core/domain/constants/injection-token.constant';
import type { ProductDetails } from '@/core/domain/entities/product.entity';
import { FingerprintRotatorService } from '@/infrastructure/proxy/fingerprint-rotator.service';
import { randomDelay } from '@/core/shared/utils/delay.util';

const NEXT_DATA_REGEX = /<script[^>]*id=["']__NEXT_DATA__["'][^>]*>([\s\S]*?)<\/script>/;
const OG_TAG_REGEX = /<meta\s+(?:[^>]*?)(?:property|name)=["']og:(\w+)["'](?:[^>]*?)content=["']([^"']*)["'][^>]*>/gi;
const OG_TAG_REVERSE_REGEX = /<meta\s+(?:[^>]*?)content=["']([^"']*)["'](?:[^>]*?)(?:property|name)=["']og:(\w+)["'][^>]*>/gi;
const TITLE_REGEX = /<title[^>]*>([\s\S]*?)<\/title>/i;
const META_DESC_REGEX = /<meta\s+(?:[^>]*?)name=["']description["'](?:[^>]*?)content=["']([^"']*)["'][^>]*>/i;
const META_DESC_REVERSE_REGEX = /<meta\s+(?:[^>]*?)content=["']([^"']*)["'](?:[^>]*?)name=["']description["'][^>]*>/i;
const META_TAG_REGEX = /<meta\s+(?:[^>]*?)(?:name|property)=["']([^"']+)["'](?:[^>]*?)content=["']([^"']*)["'][^>]*>/gi;
const META_TAG_REVERSE_REGEX = /<meta\s+(?:[^>]*?)content=["']([^"']*)["'](?:[^>]*?)(?:name|property)=["']([^"']+)["'][^>]*>/gi;
const JSON_LD_REGEX = /<script[^>]*type=["']application\/ld\+json["'][^>]*>([\s\S]*?)<\/script>/gi;

@Injectable()
export class GenericNaverScraperStrategy {
  private readonly logger = new Logger(GenericNaverScraperStrategy.name);
  private readonly limiter: Bottleneck;
  private readonly maxRetries: number;
  private readonly minDelayMs: number;
  private readonly maxDelayMs: number;

  constructor(
    @Inject(HTTP_CLIENT_PORT)
    private readonly httpClient: HttpClientPort,
    @Inject(PROXY_MANAGER_PORT)
    private readonly proxyManager: ProxyManagerPort,
    private readonly fingerprintRotator: FingerprintRotatorService,
    private readonly configService: ConfigService,
  ) {
    this.maxRetries = this.configService.get<number>('scraping.maxRetries', 3);
    this.minDelayMs = this.configService.get<number>('scraping.minDelayMs', 500);
    this.maxDelayMs = this.configService.get<number>('scraping.maxDelayMs', 2000);
    this.limiter = new Bottleneck({
      maxConcurrent: this.configService.get<number>('scraping.maxConcurrent', 5),
      minTime: this.configService.get<number>('scraping.throttleMinTimeMs', 200),
    });
  }

  async scrapePage(url: string): Promise<ProductDetails> {
    return this.limiter.schedule(() => this.scrapeWithRetry(url));
  }

  private async scrapeWithRetry(
    url: string,
    attempt = 1,
  ): Promise<ProductDetails> {
    const proxy = this.proxyManager.getNextProxy();

    try {
      await randomDelay(this.minDelayMs, this.maxDelayMs);
      const fingerprint = this.fingerprintRotator.generate();

      const response = await this.httpClient.get<string>(url, {
        headers: {
          ...fingerprint,
          Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
          Referer: 'https://www.naver.com/',
        },
        proxy: proxy?.url,
        waitUntil: 'networkidle',
        waitForSelector: 'title, meta[property="og:title"], h1, [class*="title"]',
      });

      const html = typeof response.data === 'string'
        ? response.data
        : JSON.stringify(response.data);

      return this.parseHtml(html, url);
    } catch (error) {
      const isRateLimited = (error as NodeJS.ErrnoException).code === 'ERR_RATE_LIMITED';
      if (proxy && !isRateLimited) {
        this.proxyManager.markProxyFailed(proxy);
      }
      if (attempt < this.maxRetries) {
        const backoff = isRateLimited ? 3000 * attempt : this.minDelayMs;
        this.logger.warn(`Retry ${attempt}/${this.maxRetries} for ${url} (wait ${backoff}ms)`);
        await randomDelay(backoff, backoff + 1000);
        return this.scrapeWithRetry(url, attempt + 1);
      }
      throw error;
    }
  }

  private parseHtml(html: string, url: string): ProductDetails {
    const result: ProductDetails = {
      url,
      source: 'generic-naver-html',
    };

    const ogTags = this.extractOgTags(html);

    result.title = ogTags['title']
      ?? this.extractTitle(html);

    result.description = ogTags['description']
      ?? this.extractMetaDescription(html);

    if (ogTags['image']) {
      result.image = ogTags['image'];
    }

    if (Object.keys(ogTags).length > 0) {
      result.ogTags = ogTags;
    }

    const metaTags = this.extractAllMetaTags(html);
    if (Object.keys(metaTags).length > 0) {
      result.metaTags = metaTags;
    }

    const nextDataMatch = html.match(NEXT_DATA_REGEX);
    if (nextDataMatch) {
      try {
        const nextData = JSON.parse(nextDataMatch[1].trim());
        result.nextData = nextData;
        this.enrichFromNextData(result, nextData);
      } catch {
        this.logger.debug('Could not parse __NEXT_DATA__');
      }
    }

    const jsonLdScripts: unknown[] = [];
    let jsonLdMatch;
    while ((jsonLdMatch = JSON_LD_REGEX.exec(html)) !== null) {
      try {
        jsonLdScripts.push(JSON.parse(jsonLdMatch[1].trim()));
      } catch {
        this.logger.debug('Could not parse JSON-LD script');
      }
    }
    if (jsonLdScripts.length > 0) {
      result.jsonLd = jsonLdScripts;
      this.enrichFromJsonLd(result, jsonLdScripts);
    }

    return result;
  }

  private enrichFromNextData(result: ProductDetails, nextData: Record<string, unknown>): void {
    const props = nextData['props'] as Record<string, unknown> | undefined;
    const pageProps = props?.['pageProps'] as Record<string, unknown> | undefined;
    if (!pageProps) return;

    if (!result.title) {
      result.title = this.findDeepValue(pageProps, ['title', 'name', 'productName', 'fundingTitle']);
    }
    if (!result.description) {
      result.description = this.findDeepValue(pageProps, ['description', 'summary', 'content']);
    }
    if (!result.image) {
      result.image = this.findDeepValue(pageProps, ['image', 'imageUrl', 'thumbnailUrl', 'thumbnail', 'mainImage']);
    }
    if (!result.price) {
      result.price = this.findDeepValue(pageProps, ['price', 'salePrice', 'currentAmount', 'amount']);
    }
  }

  private enrichFromJsonLd(result: ProductDetails, jsonLdScripts: unknown[]): void {
    for (const ld of jsonLdScripts) {
      if (!ld || typeof ld !== 'object') continue;
      const ldObj = ld as Record<string, unknown>;

      if (!result.title && ldObj['name']) {
        result.title = ldObj['name'] as string;
      }
      if (!result.description && ldObj['description']) {
        result.description = ldObj['description'] as string;
      }
      if (!result.image) {
        const img = ldObj['image'];
        if (typeof img === 'string') {
          result.image = img;
        } else if (Array.isArray(img) && img.length > 0) {
          const first = img[0];
          result.image = typeof first === 'string' ? first : String((first as Record<string, unknown>)?.['url'] ?? '');
        } else if (img && typeof img === 'object') {
          result.image = String((img as Record<string, unknown>)['url'] ?? '');
        }
      }
    }
  }

  private findDeepValue(obj: Record<string, unknown>, keys: string[]): string | undefined {
    for (const key of keys) {
      const val = this.deepGet(obj, key);
      if (val !== undefined && val !== null && val !== '') {
        return String(val);
      }
    }
    return undefined;
  }

  private deepGet(obj: unknown, targetKey: string, depth = 0): unknown {
    if (depth > 5 || !obj || typeof obj !== 'object') return undefined;

    const record = obj as Record<string, unknown>;
    if (record[targetKey] !== undefined && record[targetKey] !== null) {
      return record[targetKey];
    }

    for (const val of Object.values(record)) {
      if (val && typeof val === 'object') {
        const found = this.deepGet(val, targetKey, depth + 1);
        if (found !== undefined && found !== null) return found;
      }
    }
    return undefined;
  }

  private extractOgTags(html: string): Record<string, string> {
    const tags: Record<string, string> = {};

    let match;
    while ((match = OG_TAG_REGEX.exec(html)) !== null) {
      tags[match[1]] = this.decodeHtmlEntities(match[2]);
    }

    while ((match = OG_TAG_REVERSE_REGEX.exec(html)) !== null) {
      const key = match[2];
      if (!tags[key]) {
        tags[key] = this.decodeHtmlEntities(match[1]);
      }
    }

    return tags;
  }

  private extractTitle(html: string): string | undefined {
    const match = html.match(TITLE_REGEX);
    if (!match) return undefined;
    const title = match[1].trim();
    return title ? this.decodeHtmlEntities(title) : undefined;
  }

  private extractMetaDescription(html: string): string | undefined {
    const match = html.match(META_DESC_REGEX) ?? html.match(META_DESC_REVERSE_REGEX);
    if (!match) return undefined;
    return this.decodeHtmlEntities(match[1]);
  }

  private extractAllMetaTags(html: string): Record<string, string> {
    const tags: Record<string, string> = {};

    let match;
    while ((match = META_TAG_REGEX.exec(html)) !== null) {
      const key = match[1];
      if (!key.startsWith('og:')) {
        tags[key] = this.decodeHtmlEntities(match[2]);
      }
    }

    while ((match = META_TAG_REVERSE_REGEX.exec(html)) !== null) {
      const key = match[2];
      if (!key.startsWith('og:') && !tags[key]) {
        tags[key] = this.decodeHtmlEntities(match[1]);
      }
    }

    return tags;
  }

  private decodeHtmlEntities(text: string): string {
    return text
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&quot;/g, '"')
      .replace(/&#39;/g, "'");
  }
}
