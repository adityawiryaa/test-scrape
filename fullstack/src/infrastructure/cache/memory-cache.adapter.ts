import { Injectable } from '@nestjs/common';
import type { CachePort } from '@/core/application/ports/cache.port';

interface CacheEntry<T> {
  value: T;
  expiresAt: number;
}

@Injectable()
export class MemoryCacheAdapter implements CachePort {
  private readonly store = new Map<string, CacheEntry<unknown>>();

  async get<T = unknown>(key: string): Promise<T | null> {
    const entry = this.store.get(key) as CacheEntry<T> | undefined;
    if (!entry) return null;
    if (Date.now() > entry.expiresAt) {
      this.store.delete(key);
      return null;
    }
    return entry.value;
  }

  async set<T = unknown>(
    key: string,
    value: T,
    ttlSeconds: number,
  ): Promise<void> {
    this.store.set(key, {
      value,
      expiresAt: Date.now() + ttlSeconds * 1000,
    });
  }

  async del(key: string): Promise<void> {
    this.store.delete(key);
  }
}
