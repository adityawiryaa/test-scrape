import {
  Injectable,
  type NestInterceptor,
  type ExecutionContext,
  type CallHandler,
  Logger,
} from '@nestjs/common';
import { Observable, map, tap } from 'rxjs';
import { ApiResponse, type ApiResponseBody } from '@/core/shared/response/api-response';

@Injectable()
export class ResponseTransformInterceptor implements NestInterceptor {
  private readonly logger = new Logger(ResponseTransformInterceptor.name);

  intercept(context: ExecutionContext, next: CallHandler): Observable<ApiResponseBody> {
    const now = Date.now();
    const request = context.switchToHttp().getRequest<{ url: string }>();

    return next.handle().pipe(
      tap(() => {
        const elapsed = Date.now() - now;
        this.logger.log(`${request.url} - ${elapsed}ms`);
      }),
      map((data) => ApiResponse.success(data)),
    );
  }
}
