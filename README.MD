<div style="text-align: center;">
  <h1>
    Create your Node application in seconds
  </h1>
</div>

Node-builder is a streamlined tool designed to help you set up a Node.js application effortlessly. It integrates seamlessly with your preferred backend tools such as Fastify and Drizzle, and supports popular databases like MySQL and PostgreSQL.

*Inspired by [Go-Blueprint](https://github.com/Melkeydev/go-blueprint)*

## Features

- Quick setup of Node.js applications
- Support for multiple backend frameworks
- Integration with popular databases
- Choice of well-known ORMs

## Table of Contents

- [Installation](#Installation)
- [Frameworks Supported](#frameworks-supported)
- [Database Support](#database-support)
- [Supported ORMs](#supported-orms)
- [License](#license)

<a id="Installation"></a>

## Installation


Ensure you have Go installed on your system, then run:

```bash
go install github.com/kutaui/node-builder@latest
```

Then in a new terminal run:

```bash
node-builder create
```

<a id="frameworks-supported"></a>

## Frameworks Supported

- [Express](https://expressjs.com)
- [Fastify](https://fastify.dev)
- [Hono](https://hono.dev/) (soon)
- [Koa](https://koajs.com/) (soon)

<a id="database-support"></a>

## Database Support

- [Mysql](https://github.com/sidorares/node-mysql2)
- [Postgres](https://github.com/brianc/node-postgres)
- [Sqlite](https://github.com/TryGhost/node-sqlite3)
- [Mongo](https://github.com/mongodb/node-mongodb-native) (soon)
- [Redis](https://github.com/redis/node-redis) (soon)

<a id="supported-orms"></a>

## Supported ORMs

- [Drizzle](https://orm.drizzle.team/)
- [Prisma](https://www.prisma.io/)
- [Sequelize](https://sequelize.org/)
- [Typeorm](https://typeorm.io/) 





<a id="license"></a>

## License
Licensed under [MIT License](./LICENSE)