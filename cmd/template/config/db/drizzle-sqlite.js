import { drizzle } from 'drizzle-orm/better-sqlite3';
import Database from 'better-sqlite3';

export async function createDbConnection() {
    const sqlite = new Database('sqlite.db');
    const db = drizzle(sqlite);
    return db;
}
