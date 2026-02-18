export interface ChannelResolverPort {
  resolveChannelUid(productUrl: string): Promise<string>;
}
