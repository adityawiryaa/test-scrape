export type ProductDetails = Record<string, unknown>;

export interface ScrapedProduct {
  productId: string;
  channelUid: string;
  details: ProductDetails;
  benefits: Record<string, unknown>;
  scrapedAt: string;
}
