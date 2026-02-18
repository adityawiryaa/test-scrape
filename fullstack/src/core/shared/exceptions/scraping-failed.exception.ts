import { HttpException, HttpStatus } from '@nestjs/common';

export class ScrapingFailedException extends HttpException {
  constructor(message = 'Failed to scrape product details') {
    super(
      { code: 'SCRAPING_FAILED', message },
      HttpStatus.INTERNAL_SERVER_ERROR,
    );
  }
}
