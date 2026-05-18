# aru-backend

Backend service for the Aru (BUMN) application. Built with Go + GoFiber + PostgreSQL + MinIO, following a clean architecture pattern.

## Architecture

![Clean Architecture](architecture.png)

Request flow:

1. External system perform request (HTTP)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system 
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

## Tech Stack

- **Golang** 1.25 — https://github.com/golang/go
- **PostgreSQL** — https://github.com/postgres/postgres
- **MinIO** — https://min.io

## Framework & Library

- **GoFiber** (HTTP framework) — https://github.com/gofiber/fiber
- **GORM** (ORM) — https://github.com/go-gorm/gorm
- **golang-migrate** (DB migration) — https://github.com/golang-migrate/migrate
- **Viper** (config loader) — https://github.com/spf13/viper
- **Logrus** (logger) — https://github.com/sirupsen/logrus
- **go-playground/validator** (request validation) — https://github.com/go-playground/validator
- **golang-jwt** (JWT auth) — https://github.com/golang-jwt/jwt
- **minio-go** (MinIO SDK) — https://github.com/minio/minio-go
- **bregydoc/gtranslate** (auto-translate ID→EN for admin content) — https://github.com/bregydoc/gtranslate

## Prerequisites

- Go 1.25+
- PostgreSQL running and accessible
- MinIO running with the bucket already created (see `config.json` → `minio.bucket`)
- `migrate` CLI installed:

  ```powershell
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

  Make sure `$env:GOPATH\bin` (default `C:\Users\<user>\go\bin`) is on your `PATH`.

## Configuration

All configuration lives in `config.json` (project root). Key fields:

| Key | Default | Description |
|---|---|---|
| `web.port` | `3001` | HTTP server port |
| `web.cors.allowOrigins` | `localhost:3000` | Frontend origin |
| `database.*` | — | Postgres credentials (username, password, host, port, name) |
| `minio.endpoint` | `127.0.0.1:9000` | MinIO server |
| `minio.bucket` | `aru-assets` | Bucket for all uploaded assets |
| `jwt.secret` | — | JWT signing secret |
| `jwt.expiryHours` | `24` | JWT lifetime |

> ⚠️ `config.json` is committed. For non-`local` environments, override via env vars or a separate file — **do not** store production secrets in the repo.

## Database Migration

All migration SQL lives in `db/migrations/`, using the `golang-migrate` format (`<timestamp>_<name>.up.sql` / `.down.sql`).

### Connection String

Format: `postgres://<user>:<pass>@<host>:<port>/<dbname>?sslmode=disable`

Based on the local `config.json`:

```
postgres://postgres:erika@localhost:15552/aru-db?sslmode=disable
```

Set it once per session so you don't have to retype it:

```powershell
$env:DB_URL = "postgres://postgres:erika@localhost:15552/aru-db?sslmode=disable"
```

### Create Migration

```powershell
migrate create -ext sql -dir db/migrations <migration_name>
```

Example:

```powershell
migrate create -ext sql -dir db/migrations add_histories_is_machine_fallback
```

This creates two empty files (`.up.sql` & `.down.sql`) in `db/migrations/`. Fill in the DDL yourself.

### Run Migration

Apply all pending migrations:

```powershell
migrate -path db/migrations -database $env:DB_URL up
```

Apply N steps forward:

```powershell
migrate -path db/migrations -database $env:DB_URL up 1
```

### Rollback

Roll back the last N steps:

```powershell
migrate -path db/migrations -database $env:DB_URL down 1
```

Roll back **everything** (careful — drops the schema):

```powershell
migrate -path db/migrations -database $env:DB_URL down -all
```

### Status & Recovery

Check the current version:

```powershell
migrate -path db/migrations -database $env:DB_URL version
```

If a migration is marked "dirty" (failed mid-run), force the state to a specific version (check the `schema_migrations` table first):

```powershell
migrate -path db/migrations -database $env:DB_URL force <version>
```

## Run Application

### Run web server

```powershell
go run cmd/web/main.go
```

Defaults to listening on `http://localhost:3001`

### Run unit test

```powershell
go test -v ./test/
```

### Build binary

```powershell
go build -o app.exe cmd/web/main.go
```

Then run:

```powershell
.\app.exe
```

## API Spec

Endpoint specifications live in the `api/` folder.

## Project Structure

```
aru-backend/
├── api/                  # API spec / docs
├── cmd/web/              # Entry point (main.go)
├── config.json           # App config
├── db/migrations/        # SQL migration files
├── internal/
│   ├── config/           # Viper config loader + DI wiring
│   ├── delivery/         # HTTP handlers (Fiber routes)
│   ├── entity/           # DB entities (GORM models)
│   ├── model/            # Request/response DTOs
│   ├── repository/       # DB access layer
│   └── usecase/          # Business logic
└── test/                 # Unit tests
```