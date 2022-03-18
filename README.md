# Digital Account

Service responsible to deal with accounts and transactions.

## Starting

This is a quick guide on how to install and run the project, just follow the next steps.

### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Golang v1.17](https://golang.org)

### Installation

1. Clone the repository

```shell
git clone git@github.com:brunomdev/digital-account.git
```

2. Access the directory created after the clone, copy the file `.env.example`  and rename to `.env`.

```shell
cd digital-account
cp .env.example .env
```

> **Information**
>
> If necessary change the values of `.env`.

### Executing

To execute the project follow the commands below:

Using docker
```shell
docker-compose up -d
# OR
make docker-up
```

Locally
```shell
go mod download
go build -o main .
docker-compose up -d db.digital-account.dev
./main
# OR
make build-and-run
```

## Documentation

With the project running access the url http://localhost:8080/docs to check the API documentation.

> **Information**
>
> Depending of the values on your .env file maybe you need to change the PORT in the url or the url itself.

## Executing tests

To execute the tests use the following command:

```shell
go test ./...
```

## Main dependencies

- [Fiber v2](https://gofiber.io)
- [MariaDB](https://mariadb.org/) (Docker)
- [go-migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate)
- [Viper](https://github.com/spf13/viper)
- [Validator](https://github.com/go-playground/validator)
