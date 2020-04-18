# robotstxt-datastore ![CI](https://github.com/paulvollmer/robotstxt-datastore/workflows/CI/badge.svg) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/paulvollmer/robotstxt-datastore?style=plastic)

`robotstxt-datastore` is a grpc service to store robots.txt data at a postgres database

## Environment Variables

At the table below yo find all environment variables to configure the gRPC server 
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
| `SENTRY_DSN`             | a sentry dsn url                                             | ``                              |

## Development

clone the repository

```
git clone git@github.com:paulvollmer/robotstxt-datastore.git
```

### gRPC Server

If you change the protobuf sources, you need to recompile by running

```
cd proto
make
```

If you change the `ent/schema` you need to generate the database code

```
cd server
make ent
```

Run the tests and build the grpc server  

```
cd server
make test
make build
```

### gRPC Client

With the client you can check, list or get robots.txt data straight from your terminal by sending requests to the grpc server.

```
cd client
go build
./client check google.com
```

To refresh a robots.txt you can use the `-r` flag

```
./client check -r google.com
```
