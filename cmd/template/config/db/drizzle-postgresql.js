import { drizzle } from "drizzle-orm/node-postgres";
import { Client } from "pg";

export async function createDbConnection() {
  const client = new Client({
    connectionString: "postgres://user:password@host:port/db",
  });

  await client.connect();
  const db = drizzle(client);

  return db;
}
