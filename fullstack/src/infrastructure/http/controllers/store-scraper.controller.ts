import { Controller, Get, Query, UsePipes, ValidationPipe } from '@nestjs/common';
import { ApiOperation, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';
import { ScrapeProductService } from '@/core/application/usecases/scrape-product.service';
import { ScrapeProductRequest } from '@/core/application/dto/request/scrape-product.request';
import type { ScrapedProduct } from '@/core/domain/entities/product.entity';

@ApiTags('Naver Scraper')
@Controller('naver')
export class StoreScraperController {
  constructor(private readonly scrapeProductService: ScrapeProductService) {}

  @Get()
  @ApiOperation({ summary: 'Scrape any Naver URL' })
  @ApiQuery({
    name: 'productUrl',
    required: true,
    description: 'Any Naver URL (*.naver.com)',
    example: 'https://smartstore.naver.com/paparecipe/products/5738498489',
  })
  @ApiResponse({ status: 200, description: 'Page scraped successfully' })
  @ApiResponse({ status: 400, description: 'Invalid URL' })
  @ApiResponse({ status: 500, description: 'Scraping failed' })
  @UsePipes(new ValidationPipe({ transform: true, whitelist: true }))
  async scrapeProduct(
    @Query() query: ScrapeProductRequest,
  ): Promise<ScrapedProduct> {
    return this.scrapeProductService.execute({
      storeName: 'naver',
      productUrl: query.productUrl,
    });
  }
}
