import {
  Injectable,
  Logger,
  OnModuleInit,
  OnModuleDestroy,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { chromium } from 'patchright';
import type { Browser, BrowserContext, Page, Response } from 'patchright';
import { newInjectedContext } from 'fingerprint-injector';
import { FingerprintGenerator } from 'fingerprint-generator';
import type {
  HttpClientPort,
  HttpClientOptions,
  HttpClientResponse,
} from '@/core/application/ports/http-client.port';

interface PlaywrightProxyConfig {
  server: string;
  username?: string;
  password?: string;
}

interface StoredCookie {
  name: string;
  value: string;
  domain: string;
  path: string;
}

const NAVER_DOMAINS = [
  '.naver.com',
  '.smartstore.naver.com',
  '.brand.naver.com',
  '.search.shopping.naver.com',
  '.happybean.naver.com',
];

@Injectable()
export class PlaywrightHttpClient
  implements HttpClientPort, OnModuleInit, OnModuleDestroy
{
  private readonly logger = new Logger(PlaywrightHttpClient.name);
  private readonly defaultTimeout: number;
  private readonly fingerprintGenerator: FingerprintGenerator;
  private readonly naverBypassCookie: string;
  private readonly cookieStore: Map<string, StoredCookie[]> = new Map();
  private browser!: Browser;

  constructor(private readonly configService: ConfigService) {
    this.defaultTimeout = this.configService.get<number>(
      'scraping.timeoutMs',
      30000,
    );
    this.naverBypassCookie = this.configService.get<string>(
      'scraping.naverBypassCookie',
      '',
    );
    this.fingerprintGenerator = new FingerprintGenerator();
  }

  async onModuleInit(): Promise<void> {
    this.browser = await chromium.launch({
      headless: true,
      args: [
        '--disable-blink-features=AutomationControlled',
        '--disable-features=IsolateOrigins,site-per-process',
        '--disable-dev-shm-usage',
        '--no-first-run',
        '--no-default-browser-check',
      ],
    });
    this.logger.log('Patchright Chromium browser launched');
  }

  async onModuleDestroy(): Promise<void> {
    if (this.browser) {
      await this.browser.close();
      this.logger.log('Patchright Chromium browser closed');
    }
  }

  async get<T = unknown>(
    url: string,
    options?: HttpClientOptions,
  ): Promise<HttpClientResponse<T>> {
    const proxy = this.parseProxy(options?.proxy);
    const timeout = options?.timeout ?? this.defaultTimeout;

    const fingerprint = this.fingerprintGenerator.getFingerprint({
      locales: ['ko-KR', 'ko'],
      operatingSystems: ['macos', 'windows'],
      browsers: [{ name: 'chrome', minVersion: 120 }],
    });

    const contextOptions: Record<string, unknown> = {
      locale: 'ko-KR',
      timezoneId: 'Asia/Seoul',
      ignoreHTTPSErrors: true,
      viewport: fingerprint.fingerprint.screen
        ? {
            width: fingerprint.fingerprint.screen.width,
            height: fingerprint.fingerprint.screen.height,
          }
        : { width: 1920, height: 1080 },
    };

    if (proxy) {
      contextOptions['proxy'] = proxy;
    }

    let context: BrowserContext | undefined;
    let page: Page | undefined;

    try {
      const injectedContext = await newInjectedContext(this.browser as never, {
        fingerprint,
        newContextOptions: contextOptions as never,
      });
      context = injectedContext as unknown as BrowserContext;

      await this.injectCookies(context, url);

      page = await context.newPage();

      if (options?.headers) {
        await page.setExtraHTTPHeaders(options.headers);
      }

      const waitUntil = options?.waitUntil ?? 'domcontentloaded';
      this.logger.debug(`GET ${url} ${proxy ? 'via proxy' : 'direct'} (${waitUntil})`);

      const response = await page.goto(url, {
        waitUntil,
        timeout,
      });

      if (!response) {
        throw new Error(`No response received for ${url}`);
      }

      const status = response.status();
      const headers = this.extractHeaders(response);
      const contentType = headers['content-type'] ?? '';

      await this.storeCookiesFromContext(context, url);

      if (status === 429 || status === 490) {
        this.logger.warn(`${status} rate limited: ${url}`);
        const error = new Error(`Rate limited by server (${status})`);
        (error as NodeJS.ErrnoException).code = 'ERR_RATE_LIMITED';
        throw error;
      }

      if (options?.waitForSelector) {
        await page.waitForSelector(options.waitForSelector, { timeout: timeout / 2 }).catch(() => {
          this.logger.debug(`Selector "${options.waitForSelector}" not found, continuing`);
        });
      }

      const body = await page.content();

      let data: T;
      if (contentType.includes('application/json')) {
        const jsonText = await page.evaluate(
          () => document.body?.innerText ?? '',
        );
        data = JSON.parse(jsonText) as T;
      } else {
        data = body as T;
      }

      return { data, status, headers };
    } finally {
      if (page) await page.close().catch(() => {});
      if (context) await context.close().catch(() => {});
    }
  }

  private async injectCookies(
    context: BrowserContext,
    url: string,
  ): Promise<void> {
    const cookies: Array<{
      name: string;
      value: string;
      domain: string;
      path: string;
    }> = [];

    if (this.naverBypassCookie && this.isNaverUrl(url)) {
      for (const domain of NAVER_DOMAINS) {
        cookies.push({
          name: 'X-Wtm-Cpt-Tk',
          value: this.naverBypassCookie,
          domain,
          path: '/',
        });
      }
      this.logger.debug('Injecting bypass cookie for Naver domains');
    }

    const hostname = new URL(url).hostname;
    const stored = this.cookieStore.get(hostname);
    if (stored) {
      cookies.push(...stored);
    }

    if (cookies.length > 0) {
      await context.addCookies(cookies);
    }
  }

  private async storeCookiesFromContext(
    context: BrowserContext,
    url: string,
  ): Promise<void> {
    try {
      const hostname = new URL(url).hostname;
      const contextCookies = await context.cookies();
      if (contextCookies.length > 0) {
        this.cookieStore.set(
          hostname,
          contextCookies.map((c) => ({
            name: c.name,
            value: c.value,
            domain: c.domain,
            path: c.path,
          })),
        );
      }
    } catch {
      this.logger.debug('Could not store cookies from context');
    }
  }

  private isNaverUrl(url: string): boolean {
    try {
      const hostname = new URL(url).hostname;
      return hostname.endsWith('.naver.com') || hostname === 'naver.com';
    } catch {
      return false;
    }
  }

  private parseProxy(proxyUrl?: string): PlaywrightProxyConfig | undefined {
    if (!proxyUrl) return undefined;

    try {
      const parsed = new URL(proxyUrl);
      return {
        server: `${parsed.protocol}//${parsed.hostname}:${parsed.port}`,
        username: parsed.username || undefined,
        password: parsed.password || undefined,
      };
    } catch {
      this.logger.warn(`Invalid proxy URL: ${proxyUrl}`);
      return undefined;
    }
  }

  private extractHeaders(response: Response): Record<string, string> {
    const result: Record<string, string> = {};
    const allHeaders = response.headers();
    for (const [key, val] of Object.entries(allHeaders)) {
      if (typeof val === 'string') result[key] = val;
    }
    return result;
  }
}
