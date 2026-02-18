export interface HttpClientOptions {
  headers?: Record<string, string>;
  proxy?: string;
  timeout?: number;
  waitUntil?: 'domcontentloaded' | 'networkidle' | 'load' | 'commit';
  waitForSelector?: string;
}

export interface HttpClientResponse<T = unknown> {
  data: T;
  status: number;
  headers: Record<string, string>;
}

export interface HttpClientPort {
  get<T = unknown>(
    url: string,
    options?: HttpClientOptions,
  ): Promise<HttpClientResponse<T>>;
}
