import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type {
  ProxyManagerPort,
  ProxyInfo,
} from '@/core/application/ports/proxy-manager.port';

@Injectable()
export class ProxyRotatorService implements ProxyManagerPort {
  private readonly logger = new Logger(ProxyRotatorService.name);
  private readonly proxies: ProxyInfo[] = [];
  private currentIndex = 0;
  private readonly failedProxies = new Set<string>();

  constructor(private readonly configService: ConfigService) {
    this.loadProxies();
  }

  private loadProxies(): void {
    const host = this.configService.get<string>('scraping.proxyHost', '');
    const port = this.configService.get<number>('scraping.proxyPort', 0);
    if (!host || !port) return;

    const username =
      this.configService.get<string>('scraping.proxyUsername', '') || undefined;
    const password =
      this.configService.get<string>('scraping.proxyPassword', '') || undefined;

    const auth = username && password ? `${username}:${password}@` : '';
    const url = `http://${auth}${host}:${port}`;

    this.proxies.push({ url, host, port, username, password });
    this.logger.log(`Loaded ${this.proxies.length} proxy`);
  }

  getNextProxy(): ProxyInfo | null {
    if (this.proxies.length === 0) return null;

    const active = this.proxies.filter((p) => !this.failedProxies.has(p.url));
    if (active.length === 0) {
      this.failedProxies.clear();
      const proxy = this.proxies[this.currentIndex % this.proxies.length];
      this.currentIndex++;
      return proxy;
    }

    const proxy = active[this.currentIndex % active.length];
    this.currentIndex++;
    return proxy;
  }

  markProxyFailed(proxy: ProxyInfo): void {
    this.failedProxies.add(proxy.url);
  }

  getActiveProxyCount(): number {
    return this.proxies.filter((p) => !this.failedProxies.has(p.url)).length;
  }
}
