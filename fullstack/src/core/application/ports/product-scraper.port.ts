export interface ProductScraperPort {
  scrapeProductDetail(
    channelUid: string,
    productId: string,
  ): Promise<Record<string, unknown>>;
}
