import { drizzle } from "drizzle-orm/mysql2";
import mysql from "mysql2/promise";


export async function createDbConnection() {
    const connection = await mysql.createConnection({
        host: "host",
        user: "user",
        database: "database",
    });

    const db = drizzle(connection);

    return db;
}

