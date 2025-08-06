import { SQL } from "bun";

export const sql = new SQL({
  adapter: "mysql",
  hostname: Bun.env.DB_HOST,
  port: Bun.env.DB_PORT,
  database: Bun.env.DB_NAME,
  username: Bun.env.DB_USER,
  password: Bun.env.DB_PASS,

  max: 20,
  idleTimeout: 30,
  maxLifetime: 0,
  connectionTimeout: 30,

  tls: true,
  onconnect: client => {
    console.log("Connected to database");
  },
  onclose: client => {
    console.log("Connection closed");
  },
});