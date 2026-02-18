# Naver SmartStore Scraper API

Scalable and undetectable REST API for scraping product details from Naver SmartStore.

## Quick Start

```bash
pnpm install
cp .env.example .env   # then edit .env with your credentials
pnpm build
pnpm start
```

Server runs at `http://localhost:3000`
Swagger docs at `http://localhost:3000/api/docs`

## API Usage

All endpoints are prefixed with `/api`. Responses follow a standardized format with `requestId` (UUID) for tracking.

### Scrape Product

```
GET /api/naver?productUrl=https://smartstore.naver.com/{store}/products/{product_id}
```

**Example:**

```bash
curl "http://localhost:3000/api/naver?productUrl=https://smartstore.naver.com/paparecipe/products/5738498489"
```

**Success Response (200):**

```json
{
  "requestId": "1d1c6778-c4bf-4265-917e-971dace37997",
  "status": 200,
  "code": "SUCCESS",
  "message": "OK",
  "data": {
    "productId": "5738498489",
    "channelUid": "2sWDvWwemS4mOSxUcLvSR",
    "details": { "..." : "raw Naver product JSON" },
    "benefits": { "..." : "raw Naver benefits JSON" },
    "scrapedAt": "2026-02-10T12:00:00.000Z"
  },
  "timestamp": "2026-02-10T12:00:00.000Z"
}
```

**Error Responses:**

| Status | Code | When |
|--------|------|------|
| 200 | `SUCCESS` | Product scraped successfully |
| 400 | `BAD_REQUEST` | Missing or invalid `productUrl` param |
| 400 | `INVALID_URL` | URL doesn't match Naver SmartStore format |
| 404 | `NOT_FOUND` | Endpoint not found |
| 429 | `RATE_LIMIT_EXCEEDED` | Too many requests (100 req/min) |
| 500 | `SCRAPING_FAILED` | Failed to scrape from Naver |

### Health Check

```bash
curl http://localhost:3000/api/health
```

## Response Format

Every API response uses a consistent structure:

```typescript
{
  requestId: string;   // UUID v4 for tracking
  status: number;      // HTTP status code
  code: string;        // Machine-readable code (SUCCESS, BAD_REQUEST, etc.)
  message: string;     // Human-readable message
  data: T | null;      // Response payload (null on error)
  timestamp: string;   // ISO 8601 timestamp
}
```

## Output Data

Scraped results are saved to `data/` folder as JSON files:
```
data/paparecipe_5738498489_2026-02-10T12-00-00-000Z.json
```

## Scraping Logic

The API scrapes two Naver internal endpoints per product:

1. **Product Details:** `GET /i/v2/channels/{channelUid}/products/{productId}?withWindow=false`
2. **Benefits:** `GET /benefits/by-product?productId={id}&channelUid={uid}`

The `channelUid` is resolved by fetching the store HTML page and extracting it from `window.__PRELOADED_STATE__`.

Both endpoints return raw JSON from Naver (no field mapping/transformation).

## Anti-Detection Strategies

- **IP Rotation:** Proxy support via configurable host/port/credentials in `.env`
- **Fingerprint Rotation:** Browser fingerprint generation via `fingerprint-generator` + `fingerprint-injector`
- **Playwright Fallback:** Headless browser scraping via `patchright` (Playwright fork)
- **Request Throttling:** Bottleneck limiter (5 concurrent, 200ms min interval)
- **Random Delays:** 500ms-2000ms random delay between requests
- **Retry Logic:** Auto-retry up to 3 times on failure
- **In-Memory Cache:** 10-minute TTL to reduce redundant requests

## Environment Variables

Copy `.env.example` to `.env` and fill in your credentials.

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Server port |
| `NODE_ENV` | `development` | Environment mode |
| `PROXY_HOST` | - | Proxy hostname |
| `PROXY_PORT` | - | Proxy port |
| `PROXY_USERNAME` | - | Proxy auth username |
| `PROXY_PASSWORD` | - | Proxy auth password |
| `NAVER_BASE_URL` | `https://smartstore.naver.com` | Naver SmartStore base URL |
| `NAVER_API_BASE_URL` | `https://smartstore.naver.com/i/v2` | Naver internal API base URL |
| `NAVER_BENEFITS_PATH` | `/benefits/by-product` | Naver benefits API path |
| `NAVER_CLIENT_ID` | - | Naver API client ID |
| `NAVER_CLIENT_SECRET` | - | Naver API client secret |
| `NAVER_BYPASS_COOKIE` | - | x-wtm-cpt-tk cookie for bypass |
| `SCRAPING_MAX_RETRIES` | `3` | Max retry attempts per request |
| `SCRAPING_TIMEOUT_MS` | `30000` | HTTP request timeout in ms |
| `SCRAPING_MIN_DELAY_MS` | `500` | Min random delay between requests |
| `SCRAPING_MAX_DELAY_MS` | `2000` | Max random delay between requests |
| `SCRAPING_MAX_CONCURRENT` | `5` | Max concurrent outbound requests |
| `SCRAPING_THROTTLE_MIN_TIME_MS` | `200` | Min time between scheduled requests |
| `CACHE_TTL_SECONDS` | `600` | In-memory cache duration (10 min) |
| `THROTTLER_TTL_MS` | `60000` | Rate limit window in ms |
| `THROTTLER_LIMIT` | `100` | Max requests per rate limit window |

## Expose via Ngrok

```bash
pnpm start
ngrok http 3000
```

Share the ngrok URL for remote testing:
```
GET https://<ngrok-url>/api/naver?productUrl=https://smartstore.naver.com/paparecipe/products/5738498489
```

## Scripts

| Command | Description |
|---------|------------|
| `pnpm install` | Install dependencies |
| `pnpm start:dev` | Start dev server with hot reload |
| `pnpm start` | Start production server |
| `pnpm build` | Build for production |
| `pnpm test` | Run unit tests |
| `pnpm test:cov` | Run tests with coverage |
| `pnpm lint` | Lint and fix code |

## Tech Stack

| Package | Purpose |
|---------|---------|
| NestJS + Fastify | Framework + HTTP adapter |
| TypeScript (ES2023, strict) | Type-safe development |
| Axios | HTTP client with proxy support |
| Patchright | Playwright-based headless browser scraping |
| Fingerprint Generator/Injector | Browser fingerprint rotation |
| Bottleneck | Outbound request throttling |
| @nestjs/throttler | Inbound API rate limiting (100 req/min) |
| class-validator | DTO validation |
| @nestjs/swagger | API documentation |
| @nestjs/terminus | Health checks |
| Jest | Unit testing (85% branch, 100% line coverage) |

## Architecture

Clean Architecture with NestJS + Fastify:

```
src/
  core/
    domain/          # Entities, constants
    application/     # Use cases, DTOs, ports (interfaces)
    shared/          # Exceptions, utilities, response builder
  infrastructure/
    http/            # Controllers, interceptors, filters
    scraping/        # HTTP adapters (Axios, Playwright), strategies, factory
    proxy/           # Proxy rotation, fingerprint rotation
    cache/           # In-memory cache adapter
    config/          # App configuration
  modules/           # NestJS module wiring
```
