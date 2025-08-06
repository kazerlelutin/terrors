import { migrate } from "./migrate";

migrate(Bun.env.PG_URL || "");