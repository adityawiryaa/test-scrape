import type { ScrapedProduct, ProductDetails } from '@/core/domain/entities/product.entity';

export const VALID_PRODUCT_URL =
  'https://smartstore.naver.com/rainbows9030/products/11102379008';
export const INVALID_PRODUCT_URL = 'https://google.com/invalid';
export const STORE_NAME = 'rainbows9030';
export const PRODUCT_ID = '11102379008';
export const CHANNEL_UID = '2v1EJ3Fas87nW0bkfGZ7m';

export const mockRawProductDetail: ProductDetails = {
  url: VALID_PRODUCT_URL,
  source: 'generic-naver-html',
  title: 'Test Product',
  description: 'A great product description',
  image: 'https://example.com/img1.jpg',
};

export const mockRawBenefits: Record<string, unknown> = {
  coupons: [
    {
      couponNo: 'COUP001',
      couponName: '10% Discount',
      discountAmount: 2990,
      discountType: 'PERCENT',
    },
  ],
  point: {
    pointRate: 1,
    pointAmount: 299,
  },
};

export const mockScrapedProduct: ScrapedProduct = {
  productId: PRODUCT_ID,
  channelUid: CHANNEL_UID,
  details: mockRawProductDetail,
  benefits: mockRawBenefits,
  scrapedAt: '2026-02-10T12:00:00.000Z',
};
