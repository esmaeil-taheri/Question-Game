# GameApp

Backend for a quiz game, written in Go with the [Echo](https://echo.labstack.com/) framework.

> ⚠️ This project is under active development and will be completed over time. Currently only user authentication (register, login, profile) is implemented — the game logic itself has not been added yet.

---

## Architecture

The project follows a **layered (Clean-style) architecture** so each layer has a single, well-defined responsibility, and layers depend on each other through **interfaces** rather than concrete implementations. This keeps the code testable and lets components be swapped out (e.g. changing the database).

```
main.go                      entry point; builds config, wires services, runs the server
│
├── config/                  application-wide configuration structs
│
├── entity/                  domain models (User, Game, Question, Category, Player)
│                            no dependency on the database or HTTP
│
├── delivery/httpserver/     delivery layer
│                            HTTP handlers, routing and middleware — only connects
│                            input/output to the services, holds no business logic
│
├── service/                 business-logic (use-case) layer
│   ├── userservice/         register, login, profile + validation and domain rules
│   └── authservice/         JWT token generation and validation (access / refresh)
│
├── repository/mysql/        data-access layer; the MySQL implementation of the interfaces
│   └── migrations/          database migration files
│
└── pkg/                     generic, reusable utilities
    ├── phonenumber/         phone-number validation
    └── security/password/   password hashing and comparison with bcrypt
```

### Request flow (example: Register)

```
HTTP Request
   → delivery/httpserver (bind request)
      → service/userservice (validation, password hashing, business rules)
         → repository/mysql (persist to database)
      ← service (build response / tokens)
   ← delivery (JSON response)
```

### Key principle: Dependency Inversion

The `service` layer **defines its own interfaces** such as `Repository` and `AuthGenerator`, and the `repository` and `authservice` layers implement them. The final wiring happens in `main.go`. As a result, the service has no direct dependency on MySQL or JWT.

---

## Getting Started

### Prerequisites
- Go 1.24+
- Docker (for MySQL)

### Run the database

```bash
docker compose up -d
```

### Run the application

```bash
go run main.go
```

The server starts on port `8080`.

---

## Migrations

Migrations are managed with [`sql-migrate`](https://github.com/rubenv/sql-migrate). Its configuration lives in `repository/mysql/dbconfig.yml`.

Install the tool:

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

Commands (run from the project root, using `dbconfig.yml`):

Apply all pending migrations:

```bash
sql-migrate up -config=repository/mysql/dbconfig.yml -env=production
```

Roll back the most recent migration (one step down):

```bash
sql-migrate down -config=repository/mysql/dbconfig.yml -env=production -limit=1
```

Show migration status (applied / pending):

```bash
sql-migrate status -config=repository/mysql/dbconfig.yml -env=production
```

---

## API

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET  | `/health-check`   | Service health check     | No |
| POST | `/users/register` | Register a user          | No |
| POST | `/users/login`    | Log in and get tokens    | No |
| GET  | `/users/profile`  | View profile             | JWT (`Authorization` header) |

---

## Roadmap

Planned additions:

- [ ] **Rich Error** — a richer error type carrying a code, message and metadata instead of raw `fmt.Errorf`, so the HTTP layer can map to the correct status code (`400` / `401` / `500`).
- [ ] **`validation` package** — move validation logic out of the services into a dedicated, organized layer.
- [ ] **Complete the game logic** — implement the existing entities (`Game`, `Question`, `Category`, `Player`): creating games, managing questions, player answers and scoring.
