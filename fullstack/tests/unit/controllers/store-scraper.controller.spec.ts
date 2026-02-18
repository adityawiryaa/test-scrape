import { Test, type TestingModule } from '@nestjs/testing';
import { StoreScraperController } from '@/infrastructure/http/controllers/store-scraper.controller';
import { ScrapeProductService } from '@/core/application/usecases/scrape-product.service';
import {
  VALID_PRODUCT_URL,
  mockScrapedProduct,
} from '../../fixtures/product-detail.fixture';

describe('StoreScraperController', () => {
  let controller: StoreScraperController;
  let scrapeProductService: jest.Mocked<ScrapeProductService>;

  beforeEach(async () => {
    const mockService = {
      execute: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      controllers: [StoreScraperController],
      providers: [
        {
          provide: ScrapeProductService,
          useValue: mockService,
        },
      ],
    }).compile();

    controller = module.get<StoreScraperController>(StoreScraperController);
    scrapeProductService = module.get(ScrapeProductService);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  it('should return scraped product on success', async () => {
    scrapeProductService.execute.mockResolvedValue(mockScrapedProduct);

    const result = await controller.scrapeProduct({
      productUrl: VALID_PRODUCT_URL,
    });

    expect(result).toEqual(mockScrapedProduct);
    expect(scrapeProductService.execute).toHaveBeenCalledWith({
      storeName: 'naver',
      productUrl: VALID_PRODUCT_URL,
    });
  });

  it('should propagate errors from service', async () => {
    scrapeProductService.execute.mockRejectedValue(new Error('Scrape failed'));

    await expect(
      controller.scrapeProduct({ productUrl: VALID_PRODUCT_URL }),
    ).rejects.toThrow('Scrape failed');
  });
});
