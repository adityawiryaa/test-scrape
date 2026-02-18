import { Injectable } from '@nestjs/common';
import UserAgent from 'user-agents';

export interface Fingerprint extends Record<string, string> {
  'User-Agent': string;
  'Accept-Language': string;
  Accept: string;
  'Accept-Encoding': string;
  Connection: string;
  'Cache-Control': string;
}

const ACCEPT_LANGUAGES = [
  'ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7',
  'ko-KR,ko;q=0.9',
  'ko;q=0.9,en-US;q=0.8',
  'ko-KR,ko;q=0.8,en;q=0.6',
];

@Injectable()
export class FingerprintRotatorService {
  generate(): Fingerprint {
    const ua = new UserAgent({ deviceCategory: 'desktop' });
    const langIndex = Math.floor(Math.random() * ACCEPT_LANGUAGES.length);

    return {
      'User-Agent': ua.toString(),
      'Accept-Language': ACCEPT_LANGUAGES[langIndex],
      Accept: 'application/json, text/plain, */*',
      'Accept-Encoding': 'gzip, deflate, br',
      Connection: 'keep-alive',
      'Cache-Control': 'no-cache',
    };
  }
}
