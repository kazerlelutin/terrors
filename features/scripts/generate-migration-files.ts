import { writeFileSync } from "fs";
import path from "path";
import { getMigrationFiles, migrationsPath } from "./migration.utils";

async function generateMigrationFiles() {

  const args = process.argv.slice(2);
  const name = args?.[0]?.toLowerCase().replace("--name=", "");

  if (!name) {
    console.error("Migration name is required --name=migration_name");
    process.exit(1);
  }

  const files = getMigrationFiles();

  let lastDate: number = 0;

  for (const file of files) {
    const date = file.split("_")[0];
    const dateObj = new Date(Number(date)).getTime();
    if (dateObj > lastDate) {
      lastDate = dateObj;
    }
  }

  const newDate = lastDate ? lastDate + 1 : Date.now();
  const newFileName = `${newDate}_${name}.sql`;
  const newFilePath = path.join(migrationsPath, newFileName);

  writeFileSync(newFilePath, `
    -- Migration: ${name}
    -- Date: ${new Date().toISOString()}
    -- Description: ${name}
    `);
  console.log(`Migration file created: ${newFilePath}`);
}

generateMigrationFiles();