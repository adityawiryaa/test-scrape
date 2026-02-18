import { randomUUID } from 'crypto';

export interface ApiResponseBody<T = unknown> {
  requestId: string;
  status: number;
  code: string;
  message: string;
  data: T | null;
  timestamp: string;
}

export class ApiResponse {
  static success<T>(data: T, status = 200, message = 'OK'): ApiResponseBody<T> {
    return {
      requestId: randomUUID(),
      status,
      code: 'SUCCESS',
      message,
      data,
      timestamp: new Date().toISOString(),
    };
  }

  static error(
    status: number,
    code: string,
    message: string,
  ): ApiResponseBody<null> {
    return {
      requestId: randomUUID(),
      status,
      code,
      message,
      data: null,
      timestamp: new Date().toISOString(),
    };
  }
}
