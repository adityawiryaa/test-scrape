import { HttpException, HttpStatus } from '@nestjs/common';

export class InvalidUrlException extends HttpException {
  constructor(message = 'Invalid Naver URL format. Supported: smartstore, brand, happybean') {
    super(
      { code: 'INVALID_URL', message },
      HttpStatus.BAD_REQUEST,
    );
  }
}
