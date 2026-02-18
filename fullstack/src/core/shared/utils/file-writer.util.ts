import { writeFile, mkdir } from 'fs/promises';
import { join } from 'path';

const DATA_DIR = join(process.cwd(), 'data');

export async function saveScrapedResult(
  storeName: string,
  productId: string,
  data: unknown,
): Promise<void> {
  const storeDir = join(DATA_DIR, storeName);
  await mkdir(storeDir, { recursive: true });
  const filepath = join(storeDir, `${productId}.json`);
  await writeFile(filepath, JSON.stringify(data, null, 2), 'utf-8');
}
