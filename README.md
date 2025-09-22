# Subscription REST-api service for Effective Mobile

This REST service implements functionality for aggregating data about users' online subscriptions.

## Stack

- Golang, Echo V4 (port: 8000)
- PosgreSQL (port: 5432)
- Alloy (port: 12345)
- Loki (port: 3100)
- Grafana (port: 3000)

## Dependencies

Before running, you must install:

### Required packages

- [Docker](https://www.docker.com/)
- [Docker Compose](https://github.com/docker/compose)
- [Migrate tool](https://github.com/golang-migrate/migrate)
- [Command runner](https://github.com/casey/just)

## Usage

You can find a detailed description of the commands in the [justfile](./justfile).

### Run the application

```zsh
just run -d --build
```

### Run the documentation (Swagger, port: 8080)

```zsh
just run-swagger
```
