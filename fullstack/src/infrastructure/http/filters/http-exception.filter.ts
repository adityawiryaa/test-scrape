import {
  Catch,
  type ExceptionFilter,
  type ArgumentsHost,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { ApiResponse } from '@/core/shared/response/api-response';

@Catch()
export class GlobalExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(GlobalExceptionFilter.name);

  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const reply = ctx.getResponse();

    const status =
      exception instanceof HttpException
        ? exception.getStatus()
        : HttpStatus.INTERNAL_SERVER_ERROR;

    const { code, message } = this.extractError(exception, status);

    this.logger.error(`[${code}] ${message}`, exception instanceof Error ? exception.stack : '');

    const body = ApiResponse.error(status, code, message);
    reply.code(status).send(body);
  }

  private extractError(
    exception: unknown,
    status: number,
  ): { code: string; message: string } {
    if (exception instanceof HttpException) {
      const response = exception.getResponse();
      if (typeof response === 'object' && response !== null) {
        const obj = response as Record<string, unknown>;
        if (obj.error && typeof obj.error === 'object') {
          const err = obj.error as Record<string, unknown>;
          return {
            code: String(err.code ?? this.statusToCode(status)),
            message: String(err.message ?? exception.message),
          };
        }
        if (Array.isArray(obj.message)) {
          return {
            code: this.statusToCode(status),
            message: obj.message.join('; '),
          };
        }
        return {
          code: String(obj.code ?? this.statusToCode(status)),
          message: String(obj.message ?? exception.message),
        };
      }
      return { code: this.statusToCode(status), message: exception.message };
    }

    return {
      code: 'INTERNAL_ERROR',
      message: 'An unexpected error occurred',
    };
  }

  private statusToCode(status: number): string {
    const map: Record<number, string> = {
      400: 'BAD_REQUEST',
      401: 'UNAUTHORIZED',
      403: 'FORBIDDEN',
      404: 'NOT_FOUND',
      429: 'RATE_LIMIT_EXCEEDED',
      500: 'INTERNAL_ERROR',
    };
    return map[status] ?? 'UNKNOWN_ERROR';
  }
}
