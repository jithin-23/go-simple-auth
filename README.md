# go-simple-auth

A small REST API in Go with user management, built as a learning project while moving from NestJS to Go.

**Stack:** Go, Gin, GORM, PostgreSQL, golang-migrate

## Getting started

Requirements: Go, Docker, and the [`migrate`](https://github.com/golang-migrate/migrate) CLI.

```bash
cp .env.example .env   # adjust values if needed
make db-up             # start Postgres in Docker
make migrate-up        # create tables
make run               # start the API on :8080
```

Run `make help` to see all commands.

## Endpoints

| Method | Path                | Description       |
|--------|---------------------|-------------------|
| GET    | `/health`           | API + DB status   |
| POST   | `/api/v1/users`     | Create a user     |
| GET    | `/api/v1/users`     | List users        |
| GET    | `/api/v1/users/:id` | Get a user        |
| PATCH  | `/api/v1/users/:id` | Update a user     |
| DELETE | `/api/v1/users/:id` | Soft-delete a user|

## Project structure

```
cmd/api/            entry point
internal/
  config/           env config
  database/         Postgres connection
  server/           router setup
  middleware/       logging, panic recovery
  httpx/            request binding, error responses
  apperror/         app error types
  password/         password hashing
  user/             user model, repository, service, handler
migrations/         SQL migrations
```
