import type { ProductScraperPort } from './product-scraper.port';
import type { BenefitsScraperPort } from './benefits-scraper.port';
import type { ChannelResolverPort } from './channel-resolver.port';

export interface StoreScraperStrategy {
  productScraper: ProductScraperPort;
  benefitsScraper: BenefitsScraperPort;
  channelResolver: ChannelResolverPort;
}

export interface StoreScraperFactoryPort {
  getStrategy(storeName: string): StoreScraperStrategy;
  getSupportedStores(): string[];
}
