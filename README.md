# Subscription REST-api service for Effective Mobile

This REST service implements functionality for aggregating data about users' online subscriptions.

## Stack

- Golang, Echo V4 (port: 8000)
- PosgreSQL (port: 5432)
- Alloy (port: 12345)
- Loki (port: 3100)
- Grafana (port: 3000)

## Prerequisites

Before running, you must install:

### Required packages

- [Docker](https://www.docker.com/)
- [Docker Compose](https://github.com/docker/compose)
- [Migrate tool](https://github.com/golang-migrate/migrate)
- [Command runner](https://github.com/casey/just)

Also you must set required environment variables:

### Environment variables

```env
# App
APP_PRODUCTION=True

# Server
SERVER_HOST=""
SERVER_PORT=8000

# Databsase
DATABASE_NAME=effective_mobile
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=admin
DATABASE_PASSWORD=123

#Grafana
GRAFANA_USER=admin
GRAFANA_PASSWORD=123
```

## Usage

You can find a detailed description of the commands in the [justfile](./justfile).
Zsh is being used as default shell. If you are using another shell, change shell parameter in [justfile](./justfile).

### Run the application

```zsh
just run -d --build
```

### Run the documentation (Swagger, port: 8080)

```zsh
just run-swagger
```
