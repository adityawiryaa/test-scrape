import { Module } from '@nestjs/common';
import { CACHE_PORT } from '@/core/domain/constants/injection-token.constant';
import { MemoryCacheAdapter } from '@/infrastructure/cache/memory-cache.adapter';

@Module({
  providers: [
    {
      provide: CACHE_PORT,
      useClass: MemoryCacheAdapter,
    },
  ],
  exports: [CACHE_PORT],
})
export class CacheModule {}
