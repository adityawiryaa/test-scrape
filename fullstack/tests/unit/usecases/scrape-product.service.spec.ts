import { Test, type TestingModule } from '@nestjs/testing';
import { ConfigService } from '@nestjs/config';
import { ScrapeProductService } from '@/core/application/usecases/scrape-product.service';
import type { ProductScraperPort } from '@/core/application/ports/product-scraper.port';
import type { BenefitsScraperPort } from '@/core/application/ports/benefits-scraper.port';
import type { CachePort } from '@/core/application/ports/cache.port';
import type { ChannelResolverPort } from '@/core/application/ports/channel-resolver.port';
import type { StoreScraperFactoryPort } from '@/core/application/ports/store-scraper-factory.port';
import {
  CACHE_PORT,
  STORE_SCRAPER_FACTORY_PORT,
} from '@/core/domain/constants/injection-token.constant';
import { InvalidUrlException } from '@/core/shared/exceptions/invalid-url.exception';
import { ScrapingFailedException } from '@/core/shared/exceptions/scraping-failed.exception';
import {
  VALID_PRODUCT_URL,
  INVALID_PRODUCT_URL,
  PRODUCT_ID,
  CHANNEL_UID,
  mockRawProductDetail,
  mockRawBenefits,
} from '../../fixtures/product-detail.fixture';

describe('ScrapeProductService', () => {
  let service: ScrapeProductService;
  let productScraper: jest.Mocked<ProductScraperPort>;
  let benefitsScraper: jest.Mocked<BenefitsScraperPort>;
  let cache: jest.Mocked<CachePort>;
  let channelResolver: jest.Mocked<ChannelResolverPort>;
  let storeFactory: jest.Mocked<StoreScraperFactoryPort>;

  beforeEach(async () => {
    productScraper = {
      scrapeProductDetail: jest.fn(),
    };
    benefitsScraper = {
      scrapeBenefits: jest.fn(),
    };
    channelResolver = {
      resolveChannelUid: jest.fn().mockResolvedValue(CHANNEL_UID),
    };
    cache = {
      get: jest.fn(),
      set: jest.fn(),
      del: jest.fn(),
    };
    storeFactory = {
      getStrategy: jest.fn().mockReturnValue({
        productScraper,
        benefitsScraper,
        channelResolver,
      }),
      getSupportedStores: jest.fn().mockReturnValue(['naver']),
    };

    const mockConfigService = {
      get: jest.fn((key: string, defaultValue: unknown) => defaultValue),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ScrapeProductService,
        { provide: STORE_SCRAPER_FACTORY_PORT, useValue: storeFactory },
        { provide: CACHE_PORT, useValue: cache },
        { provide: ConfigService, useValue: mockConfigService },
      ],
    }).compile();

    service = module.get<ScrapeProductService>(ScrapeProductService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('should throw InvalidUrlException for invalid URL', async () => {
    await expect(
      service.execute({ storeName: 'naver', productUrl: INVALID_PRODUCT_URL }),
    ).rejects.toThrow(InvalidUrlException);
  });

  it('should return cached result when available', async () => {
    const cached = {
      productId: PRODUCT_ID,
      channelUid: CHANNEL_UID,
      details: mockRawProductDetail,
      benefits: mockRawBenefits,
      scrapedAt: '2026-01-01T00:00:00.000Z',
    };
    cache.get.mockResolvedValue(cached);

    const result = await service.execute({
      storeName: 'naver',
      productUrl: VALID_PRODUCT_URL,
    });

    expect(result).toEqual(cached);
    expect(productScraper.scrapeProductDetail).not.toHaveBeenCalled();
  });

  it('should scrape and cache result when not cached', async () => {
    cache.get.mockResolvedValue(null);
    productScraper.scrapeProductDetail.mockResolvedValue(mockRawProductDetail);
    benefitsScraper.scrapeBenefits.mockResolvedValue(mockRawBenefits);
    cache.set.mockResolvedValue(undefined);

    const result = await service.execute({
      storeName: 'naver',
      productUrl: VALID_PRODUCT_URL,
    });

    expect(result.productId).toBe(PRODUCT_ID);
    expect(result.channelUid).toBe(CHANNEL_UID);
    expect(result.details).toEqual(mockRawProductDetail);
    expect(result.benefits).toEqual(mockRawBenefits);
    expect(result.scrapedAt).toBeDefined();

    expect(storeFactory.getStrategy).toHaveBeenCalledWith('naver');
    expect(channelResolver.resolveChannelUid).toHaveBeenCalled();
    expect(productScraper.scrapeProductDetail).toHaveBeenCalledWith(
      CHANNEL_UID,
      PRODUCT_ID,
    );
    expect(benefitsScraper.scrapeBenefits).toHaveBeenCalledWith(
      CHANNEL_UID,
      PRODUCT_ID,
    );
    expect(cache.set).toHaveBeenCalled();
  });

  it('should throw ScrapingFailedException when scraper fails with Error', async () => {
    cache.get.mockResolvedValue(null);
    productScraper.scrapeProductDetail.mockRejectedValue(
      new Error('Network error'),
    );

    await expect(
      service.execute({ storeName: 'naver', productUrl: VALID_PRODUCT_URL }),
    ).rejects.toThrow(ScrapingFailedException);
  });

  it('should throw ScrapingFailedException when scraper fails with non-Error', async () => {
    cache.get.mockResolvedValue(null);
    productScraper.scrapeProductDetail.mockRejectedValue('string error');

    await expect(
      service.execute({ storeName: 'naver', productUrl: VALID_PRODUCT_URL }),
    ).rejects.toThrow(ScrapingFailedException);
  });
});
