import { Module } from '@nestjs/common';
import { PROXY_MANAGER_PORT } from '@/core/domain/constants/injection-token.constant';
import { ProxyRotatorService } from '@/infrastructure/proxy/proxy-rotator.service';
import { FingerprintRotatorService } from '@/infrastructure/proxy/fingerprint-rotator.service';

@Module({
  providers: [
    {
      provide: PROXY_MANAGER_PORT,
      useClass: ProxyRotatorService,
    },
    FingerprintRotatorService,
  ],
  exports: [PROXY_MANAGER_PORT, FingerprintRotatorService],
})
export class ProxyModule {}
