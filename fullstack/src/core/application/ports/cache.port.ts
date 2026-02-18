export interface CachePort {
  get<T = unknown>(key: string): Promise<T | null>;
  set<T = unknown>(key: string, value: T, ttlSeconds: number): Promise<void>;
  del(key: string): Promise<void>;
}
