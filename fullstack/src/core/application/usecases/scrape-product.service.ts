import { Inject, Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { CachePort } from '@/core/application/ports/cache.port';
import type { StoreScraperFactoryPort } from '@/core/application/ports/store-scraper-factory.port';
import type { ScrapedProduct } from '@/core/domain/entities/product.entity';
import {
  CACHE_PORT,
  STORE_SCRAPER_FACTORY_PORT,
} from '@/core/domain/constants/injection-token.constant';
import { ScrapingFailedException } from '@/core/shared/exceptions/scraping-failed.exception';
import { InvalidUrlException } from '@/core/shared/exceptions/invalid-url.exception';
import { saveScrapedResult } from '@/core/shared/utils/file-writer.util';
import type { ScrapeProductInput } from '@/core/application/dto/request/scrape-product-input.dto';
import { NAVER_URL_REGEX } from '@/core/domain/constants/scraping.constant';

@Injectable()
export class ScrapeProductService {
  private readonly logger = new Logger(ScrapeProductService.name);
  private readonly cacheTtlSeconds: number;

  constructor(
    @Inject(CACHE_PORT)
    private readonly cache: CachePort,
    @Inject(STORE_SCRAPER_FACTORY_PORT)
    private readonly storeScraperFactory: StoreScraperFactoryPort,
    private readonly configService: ConfigService,
  ) {
    this.cacheTtlSeconds = this.configService.get<number>(
      'scraping.cacheTtlSeconds',
      600,
    );
  }

  async execute(input: ScrapeProductInput): Promise<ScrapedProduct> {
    if (!NAVER_URL_REGEX.test(input.productUrl)) {
      throw new InvalidUrlException(
        'productUrl must be a valid Naver SmartStore URL (https://smartstore.naver.com/{store}/products/{id})',
      );
    }

    const cacheKey = `naver:${input.productUrl}`;
    const cached = await this.cache.get<ScrapedProduct>(cacheKey);
    if (cached) {
      this.logger.log(`Cache hit for ${input.productUrl}`);
      return cached;
    }

    try {
      const productId = this.extractProductId(input.productUrl);
      const strategy = this.storeScraperFactory.getStrategy(input.storeName);

      const channelUid = await strategy.channelResolver.resolveChannelUid(
        input.productUrl,
      );

      const [details, benefits] = await Promise.all([
        strategy.productScraper.scrapeProductDetail(channelUid, productId),
        strategy.benefitsScraper.scrapeBenefits(channelUid, productId),
      ]);

      const result: ScrapedProduct = {
        productId,
        channelUid,
        details,
        benefits,
        scrapedAt: new Date().toISOString(),
      };

      await this.cache.set(cacheKey, result, this.cacheTtlSeconds);
      await saveScrapedResult(channelUid, productId, result).catch((err) =>
        this.logger.warn(`Failed to save result to file: ${err}`),
      );

      return result;
    } catch (error) {
      if (error instanceof InvalidUrlException) throw error;
      const cause = error instanceof Error ? error.message : String(error);
      this.logger.error(
        `Failed to scrape ${input.productUrl}: ${cause}`,
        error instanceof Error ? error.stack : '',
      );
      throw new ScrapingFailedException(
        `Failed to scrape ${input.productUrl}: ${cause}`,
      );
    }
  }

  private extractProductId(url: string): string {
    const match = url.match(/\/products\/(\d+)/);
    return match?.[1] ?? 'unknown';
  }
}
