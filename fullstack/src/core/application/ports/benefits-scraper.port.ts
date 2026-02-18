export interface BenefitsScraperPort {
  scrapeBenefits(
    channelUid: string,
    productId: string,
  ): Promise<Record<string, unknown>>;
}
