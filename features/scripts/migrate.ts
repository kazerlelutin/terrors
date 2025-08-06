import { SQL } from "bun";
import path from "path";
import { getMigrationFiles } from "./migration.utils";

export type MigrationConfig = {
  host: string;
  port: number;
  user: string;
  password: string;
  database: string;
};

export async function migrate(url: string) {

  const sql = new SQL({
    url
  });

  await sql`CREATE TABLE IF NOT EXISTS migration_table (id SERIAL PRIMARY KEY, name VARCHAR(255))`;

  const migrationFiles = getMigrationFiles();

  for (const file of migrationFiles) {
    const migrationName = path.basename(file, '.sql');
    const existingMigrations = await sql`SELECT * FROM migration_table WHERE name = ${migrationName}`;

    if (existingMigrations.length > 0) {
      console.log(`Migration ${migrationName} already applied`);
      continue;
    }

    console.log(`Applying migration: ${migrationName}`);

    const text = await Bun.file(file).text();
    await sql.unsafe(text);
    await sql`INSERT INTO migration_table (name) VALUES (${migrationName})`;

    console.log(`Migration ${migrationName} applied successfully`);
  }

  console.log("Migrations applied successfully");
  process.exit(0);
}