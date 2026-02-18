import { NestFactory } from '@nestjs/core';
import {
  FastifyAdapter,
  type NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { ValidationPipe, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
  );

  app.setGlobalPrefix('api');
  app.enableCors();
  app.useGlobalPipes(
    new ValidationPipe({
      transform: true,
      whitelist: true,
      forbidNonWhitelisted: true,
    }),
  );

  const swaggerConfig = new DocumentBuilder()
    .setTitle('Naver SmartStore Scraper API')
    .setDescription('Scalable API for scraping Naver SmartStore product details')
    .setVersion('1.0')
    .build();

  const document = SwaggerModule.createDocument(app, swaggerConfig);
  SwaggerModule.setup('api/docs', app, document);

  const configService = app.get(ConfigService);
  const port = configService.get<number>('app.port', 3000);

  app.enableShutdownHooks();

  const shutdown = async (signal: string) => {
    Logger.log(`Received ${signal}, shutting down gracefully...`, 'Bootstrap');
    await app.close();
    process.exit(0);
  };

  process.on('SIGTERM', () => shutdown('SIGTERM'));
  process.on('SIGINT', () => shutdown('SIGINT'));

  try {
    await app.listen(port, '0.0.0.0');
    Logger.log(`Application running on http://localhost:${port}`, 'Bootstrap');
    Logger.log(
      `Swagger docs at http://localhost:${port}/api/docs`,
      'Bootstrap',
    );
  } catch (err: unknown) {
    const error = err as NodeJS.ErrnoException;
    if (error.code === 'EADDRINUSE') {
      Logger.error(
        `Port ${port} is already in use. Kill it with: lsof -ti:${port} | xargs kill -9`,
        'Bootstrap',
      );
      process.exit(1);
    }
    throw err;
  }
}

bootstrap();
