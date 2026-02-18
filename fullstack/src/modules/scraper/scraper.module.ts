import { Module } from '@nestjs/common';
import {
  HTTP_CLIENT_PORT,
  GENERIC_SCRAPER,
} from '@/core/domain/constants/injection-token.constant';
import { ScrapeProductService } from '@/core/application/usecases/scrape-product.service';
import { PlaywrightHttpClient } from '@/infrastructure/scraping/adapters/playwright-http.client';
import { GenericNaverScraperStrategy } from '@/infrastructure/scraping/strategies/generic-naver-scraper.strategy';
import { StoreScraperController } from '@/infrastructure/http/controllers/store-scraper.controller';
import { ProxyModule } from '@/modules/proxy/proxy.module';
import { CacheModule } from '@/modules/cache/cache.module';

@Module({
  imports: [ProxyModule, CacheModule],
  controllers: [StoreScraperController],
  providers: [
    {
      provide: HTTP_CLIENT_PORT,
      useClass: PlaywrightHttpClient,
    },
    {
      provide: GENERIC_SCRAPER,
      useClass: GenericNaverScraperStrategy,
    },
    ScrapeProductService,
  ],
  exports: [ScrapeProductService],
})
export class ScraperModule {}
