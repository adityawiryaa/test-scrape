import { registerAs } from '@nestjs/config';

export const appConfig = registerAs('app', () => ({
  port: parseInt(process.env['PORT'] ?? '3000', 10),
  environment: process.env['NODE_ENV'] ?? 'development',
}));

export const scrapingConfig = registerAs('scraping', () => ({
  maxRetries: parseInt(process.env['SCRAPING_MAX_RETRIES'] ?? '3', 10),
  timeoutMs: parseInt(process.env['SCRAPING_TIMEOUT_MS'] ?? '30000', 10),
  minDelayMs: parseInt(process.env['SCRAPING_MIN_DELAY_MS'] ?? '500', 10),
  maxDelayMs: parseInt(process.env['SCRAPING_MAX_DELAY_MS'] ?? '2000', 10),
  maxConcurrent: parseInt(process.env['SCRAPING_MAX_CONCURRENT'] ?? '5', 10),
  throttleMinTimeMs: parseInt(process.env['SCRAPING_THROTTLE_MIN_TIME_MS'] ?? '200', 10),
  cacheTtlSeconds: parseInt(process.env['CACHE_TTL_SECONDS'] ?? '600', 10),
  naverBaseUrl: process.env['NAVER_BASE_URL'] ?? 'https://smartstore.naver.com',
  naverApiBaseUrl: process.env['NAVER_API_BASE_URL'] ?? 'https://smartstore.naver.com/i/v2',
  naverBenefitsPath: process.env['NAVER_BENEFITS_PATH'] ?? '/benefits/by-product',
  naverClientId: process.env['NAVER_CLIENT_ID'] ?? '',
  naverClientSecret: process.env['NAVER_CLIENT_SECRET'] ?? '',
  naverBypassCookie: process.env['NAVER_BYPASS_COOKIE'] ?? '',
  proxyHost: process.env['PROXY_HOST'] ?? '',
  proxyPort: parseInt(process.env['PROXY_PORT'] ?? '0', 10),
  proxyUsername: process.env['PROXY_USERNAME'] ?? '',
  proxyPassword: process.env['PROXY_PASSWORD'] ?? '',
}));

export const throttlerConfig = registerAs('throttler', () => ({
  ttl: parseInt(process.env['THROTTLER_TTL_MS'] ?? '60000', 10),
  limit: parseInt(process.env['THROTTLER_LIMIT'] ?? '100', 10),
}));
