# Naver SmartStore Scraper API

Scalable and undetectable REST API for scraping product details from Naver SmartStore.

## Quick Start

```bash
pnpm install
cp .env.example .env   # then edit .env with your proxy credentials
pnpm build
pnpm start
```

Server runs at `http://localhost:3000`
Swagger docs at `http://localhost:3000/api/docs`

## API Usage

All responses follow a standardized format with `requestId` (UUID) for tracking.

### Scrape Product

```
GET /naver?productUrl=https://smartstore.naver.com/{store}/products/{product_id}
```

**Example:**

```bash
curl "http://localhost:3000/naver?productUrl=https://smartstore.naver.com/paparecipe/products/5738498489"
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
    "details": { "raw Naver product JSON (full structure)" : "..." },
    "benefits": { "raw Naver benefits JSON (full structure)" : "..." },
    "scrapedAt": "2026-02-10T12:00:00.000Z"
  },
  "timestamp": "2026-02-10T12:00:00.000Z"
}
```

**Error Responses:**

```json
{
  "requestId": "cf140cf9-2a43-4889-9bf7-c15678976ef0",
  "status": 400,
  "code": "BAD_REQUEST",
  "message": "productUrl must be a valid Naver SmartStore URL (https://smartstore.naver.com/{store}/products/{id})",
  "data": null,
  "timestamp": "2026-02-10T06:19:48.356Z"
}
```

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
curl http://localhost:3000/health
```

```json
{
  "requestId": "5ef90790-e735-4bfd-b045-48d76821c5e5",
  "status": 200,
  "code": "SUCCESS",
  "message": "OK",
  "data": { "status": "ok" },
  "timestamp": "2026-02-10T06:19:53.964Z"
}
```

## Response Format

Every API response uses a consistent structure:

```typescript
{
  requestId: string;   // UUID v4 for tracking
  status: number;      // HTTP status code
  code: string;        // Machine-readable code (SUCCESS, BAD_REQUEST, SCRAPING_FAILED, etc.)
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

- **IP Rotation:** Proxy support via configurable `PROXY_URL` in `.env`
- **Fingerprint Rotation:** Random User-Agent, Accept-Language, and browser headers per request
- **Request Throttling:** Bottleneck limiter (5 concurrent, 200ms min interval)
- **Random Delays:** 500ms-2000ms random delay between requests
- **Retry Logic:** Auto-retry up to 3 times on failure
- **In-Memory Cache:** 10-minute TTL to reduce redundant requests

## Environment Variables

All configuration is driven by `.env` — no hardcoded values in source code. Copy `.env.example` to `.env` and fill in your credentials.

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Server port |
| `NODE_ENV` | `development` | Environment mode |
| `PROXY_URL` | - | Proxy in `host:port:user:pass` format |
| `NAVER_BASE_URL` | `https://smartstore.naver.com` | Naver SmartStore base URL |
| `NAVER_API_BASE_URL` | `https://smartstore.naver.com/i/v2` | Naver internal API base URL |
| `NAVER_BENEFITS_PATH` | `/benefits/by-product` | Naver benefits API path |
| `SCRAPING_MAX_RETRIES` | `3` | Max retry attempts per request |
| `SCRAPING_TIMEOUT_MS` | `30000` | HTTP request timeout in ms |
| `SCRAPING_MIN_DELAY_MS` | `500` | Min random delay between requests |
| `SCRAPING_MAX_DELAY_MS` | `2000` | Max random delay between requests |
| `SCRAPING_MAX_CONCURRENT` | `5` | Max concurrent outbound requests |
| `SCRAPING_THROTTLE_MIN_TIME_MS` | `200` | Min time between scheduled requests |
| `CACHE_TTL_SECONDS` | `600` | In-memory cache duration (10 min) |
| `THROTTLER_TTL_MS` | `60000` | Rate limit window in ms |
| `THROTTLER_LIMIT` | `100` | Max requests per rate limit window |

If `PROXY_URL` is set, all requests go through the proxy. Otherwise, requests are made directly.

## Expose via Ngrok

```bash
pnpm start
ngrok http 3000
```

Share the ngrok URL for remote testing:
```
GET https://<ngrok-url>/naver?productUrl=https://smartstore.naver.com/paparecipe/products/5738498489
```

## Scripts

| Command | Description |
|---------|------------|
| `pnpm install` | Install dependencies |
| `pnpm start:dev` | Start dev server with hot reload |
| `pnpm start` | Start production server |
| `pnpm build` | Build for production |
| `pnpm test` | Run unit tests |
| `pnpm lint` | Lint and fix code |

## Tech Stack

| Package | Purpose |
|---------|---------|
| NestJS + Fastify | Framework + HTTP adapter |
| TypeScript | Type-safe development |
| Axios | HTTP client with proxy support |
| Bottleneck | Outbound request throttling |
| @nestjs/throttler | Inbound API rate limiting (100 req/min) |
| user-agents | Browser fingerprint generation |
| class-validator | DTO validation |
| @nestjs/swagger | API documentation |
| @nestjs/terminus | Health checks |

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
    scraping/        # HTTP adapters, scraping strategies, factory
    proxy/           # Proxy rotation, fingerprint rotation
    cache/           # In-memory cache adapter
    config/          # App configuration
  modules/           # NestJS module wiring
```
