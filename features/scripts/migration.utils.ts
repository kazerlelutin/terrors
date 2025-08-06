import { existsSync, mkdirSync, readdirSync } from 'fs';
import path from 'path';

export const migrationsPath = path.join(process.cwd(), "migrations");
export function getMigrationFiles() {
  const isExist = existsSync(migrationsPath);
  if (!isExist) mkdirSync(migrationsPath);

  const files = readdirSync(migrationsPath).filter(file => file.endsWith(".sql"));
  return files.map(file => path.join(migrationsPath, file)).sort((a, b) => {
    const dateA = a.split("_")[0];
    const dateB = b.split("_")[0];
    return Number(dateA) - Number(dateB);
  });
}