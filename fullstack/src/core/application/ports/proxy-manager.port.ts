export interface ProxyInfo {
  url: string;
  host: string;
  port: number;
  username?: string;
  password?: string;
}

export interface ProxyManagerPort {
  getNextProxy(): ProxyInfo | null;
  markProxyFailed(proxy: ProxyInfo): void;
  getActiveProxyCount(): number;
}
