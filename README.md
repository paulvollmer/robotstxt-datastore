# robotstxt-datastore ![CI](https://github.com/paulvollmer/robotstxt-datastore/workflows/CI/badge.svg) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/paulvollmer/robotstxt-datastore?style=plastic)

`robotstxt-datastore` is a gRPC service to store robots.txt data at a postgres database.

## Deploy

You can use the following `docker-compose` snippet to deploy a postgres database and the gRPC server.

```yaml
version: "3"
services:

  # the postgres database container
  postgres:
    image: postgres:10-alpine
    container_name: robotstxt_datastore_postgres
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: robotstxt
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - robotstxt_datastore_net
    volumes:
      - ./postgres-data:/var/lib/postgresql/data

# the grpc server container
  server:
    image: paulvollmer/robotstxt-datastore:v0.1.1
    container_name: robotstxt_datastoreer_server
    restart: always
    depends_on:
      - postgres
    environment:
      DATABASE_HOST: robotstxt_datastore_postgres
      DATABASE_PORT: 5432
      DATABASE_USER: postgres
      DATABASE_PASSWORD: password
      DATABASE_NAME: robotstxt
      # SENTRY_DSN: http://dsn@localhost:9100/2
    ports:
      - "5000:5000"
    networks:
      - robotstxt_datastore_net

networks:
  robotstxt_datastore_net:
```

## Environment Variables

At the table below yo find all environment variables to configure the gRPC server.

| env-var                  | description                                                  | default value                   |
| ------------------------ | ------------------------------------------------------------ | ------------------------------- |
| `SERVER_ADDR`            | the grpc server address                                      | `:5000`                         |
| `DATABASE_HOST`          | the database host                                            | `localhost`                     |
| `DATABASE_PORT`          | the database port                                            | `5432`                          |
| `DATABASE_USER`          | the database user                                            | `postgres`                      |
| `DATABASE_PASSWORD`      | the database password                                        | `password`                      |
| `DATABASE_NAME`          | the database name                                            | `robotstxt`                     |
| `DATABASE_SSLMODE`       | the database ssl mode                                        | `disable`                       |
| `REFRESH_AFTER`          | the delay after the robots.txt will be reloaded (in seconds) | `864000` default set to 10 days |
| `DEFAULT_REQUEST_SCHEME` | the default scheme used to send requests                     | `https`                         |
| `DEFAULT_LIMIT`          | the default list limit                                       | `100`                           |
| `USERAGENT`              | the User-Agent used to send the request                      | `robotstxtbot`                  |
| `SENTRY_DSN`             | a sentry dsn url                                             |                                 |

## Development

Clone the repository

```sh
git clone git@github.com:paulvollmer/robotstxt-datastore.git
```

### gRPC Server

If you change the protobuf sources, you need to recompile by running

```sh
cd proto
make
```

If you change the `ent/schema` you need to generate the database code

```sh
cd server
make ent
```

Run the tests and build the grpc server  

```sh
cd server
make test
make build
```

### gRPC Client

With the client you can check, list or get robots.txt data straight from your terminal by sending requests to the grpc server.

```sh
cd client
go build
./client check google.com
```

To refresh a robots.txt you can use the `-r` flag

```sh
./client check -r google.com
```
