import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import axios from 'axios';
import type { AxiosRequestConfig } from 'axios';
import { HttpsProxyAgent } from 'https-proxy-agent';
import type {
  HttpClientPort,
  HttpClientOptions,
  HttpClientResponse,
} from '@/core/application/ports/http-client.port';

@Injectable()
export class AxiosHttpClient implements HttpClientPort {
  private readonly logger = new Logger(AxiosHttpClient.name);
  private readonly defaultTimeout: number;

  constructor(private readonly configService: ConfigService) {
    this.defaultTimeout = this.configService.get<number>('scraping.timeoutMs', 30000);
  }

  async get<T = unknown>(
    url: string,
    options?: HttpClientOptions,
  ): Promise<HttpClientResponse<T>> {
    const config: AxiosRequestConfig = {
      timeout: options?.timeout ?? this.defaultTimeout,
      headers: options?.headers,
      proxy: false,
      validateStatus: (status) => status < 500,
    };

    if (options?.proxy) {
      const agent = new HttpsProxyAgent(options.proxy);
      config.httpsAgent = agent;
      config.httpAgent = agent;
    }

    this.logger.debug(`GET ${url} ${options?.proxy ? `via proxy` : 'direct'}`);
    const response = await axios.get<T>(url, config);

    if (response.status === 429) {
      this.logger.warn(`429 rate limited: ${url}`);
      const error = new Error('Rate limited by server (429)');
      (error as NodeJS.ErrnoException).code = 'ERR_RATE_LIMITED';
      throw error;
    }

    const headers: Record<string, string> = {};
    for (const [key, val] of Object.entries(response.headers)) {
      if (typeof val === 'string') headers[key] = val;
    }

    return {
      data: response.data,
      status: response.status,
      headers,
    };
  }
}
