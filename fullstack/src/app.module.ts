import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ThrottlerModule } from '@nestjs/throttler';
import { APP_FILTER, APP_GUARD, APP_INTERCEPTOR } from '@nestjs/core';
import { ThrottlerGuard } from '@nestjs/throttler';
import { appConfig, scrapingConfig, throttlerConfig } from '@/infrastructure/config/app.config';
import { ScraperModule } from '@/modules/scraper/scraper.module';
import { HealthModule } from '@/modules/health/health.module';
import { ResponseTransformInterceptor } from '@/infrastructure/http/interceptors/response-transform.interceptor';
import { GlobalExceptionFilter } from '@/infrastructure/http/filters/http-exception.filter';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [appConfig, scrapingConfig, throttlerConfig],
    }),
    ThrottlerModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (config: ConfigService) => [
        {
          ttl: config.get<number>('throttler.ttl', 60000),
          limit: config.get<number>('throttler.limit', 100),
        },
      ],
    }),
    ScraperModule,
    HealthModule,
  ],
  providers: [
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
    {
      provide: APP_FILTER,
      useClass: GlobalExceptionFilter,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ResponseTransformInterceptor,
    },
  ],
})
export class AppModule {}
