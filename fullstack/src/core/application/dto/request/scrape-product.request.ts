import { IsNotEmpty, IsString, Matches } from 'class-validator';
import { NAVER_URL_REGEX } from '@/core/domain/constants/scraping.constant';

export class ScrapeProductRequest {
  @IsString()
  @IsNotEmpty()
  @Matches(NAVER_URL_REGEX, {
    message:
      'productUrl must be a valid Naver SmartStore URL (https://smartstore.naver.com/{store}/products/{id})',
  })
  productUrl!: string;
}
