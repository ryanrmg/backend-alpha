# Backend Alpha

Backend Alpha is a Go REST API for interacting with trading data. The application exposes HTTP endpoints backed by PostgreSQL and is organized using a layered architecture.

## Architecture

```text
HTTP Handler
    ↓
Service
    ↓
Repository
    ↓
PostgreSQL
```

## Project Structure

```text
backend-alpha/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── api/
│   ├── config/
│   ├── database/
│   ├── repository/
│   ├── server/
│   └── service/
│
├── go.mod
├── go.sum
├── README.md
└── environment.sh
```

## Requirements

- Go 1.24+
- PostgreSQL 15+
- Git

## Installation

Clone the repository:

```bash
git clone git@github.com:ryanrmg/backend-alpha.git
cd backend-alpha
```

Load the development environment:

```bash
source environment.sh
```

Download dependencies:

```bash
go mod tidy
```

## Running the Server

```bash
go run ./cmd/server
```

or build the binary:

```bash
go build -o backend-alpha ./cmd/server
./backend-alpha
```

## Running Tests

Run all tests:

```bash
go test ./...
```

Run a specific package:

```bash
go test ./internal/service
```

Run with verbose output:

```bash
go test -v ./...
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| PORT | HTTP server port |
| DB_USER | PostgreSQL username |
| DB_PASSWORD | PostgreSQL password |
| DB_HOST | PostgreSQL hostname |
| DB_PORT | PostgreSQL port |
| DB_NAME | Development database |
| DATABASE_URL | Full PostgreSQL connection string |
| TEST_DATABASE_URL | PostgreSQL connection string used by integration tests |

## Example Request

```http
GET /api/user/trades?accountId=1234
```

Example response:

```json
[
  {
    "id": 1,
    "accountId": 1234,
    "contractId": "CON.F.US.MNQ.U26",
    "creationTimestamp": "2026-07-07T14:30:00Z",
    "price": 23567.25,
    "profitAndLoss": 125.5,
    "fees": 2.04,
    "side": 1,
    "size": 2,
    "voided": false,
    "orderId": 1001
  }
]
```

## Application Flow

```text
Browser
    ↓
Router
    ↓
TradeHandler
    ↓
TradeService
    ↓
TradeRepository (interface)
    ↓
PostgresTradeRepository
    ↓
PostgreSQL
```

## Responsibilities

| Package | Responsibility |
|---------|----------------|
| `cmd/server` | Application entry point |
| `internal/config` | Loads application configuration |
| `internal/database` | Creates and manages PostgreSQL connection pool |
| `internal/repository` | Database interfaces and PostgreSQL implementations |
| `internal/service` | Business logic |
| `internal/api` | HTTP handlers and routing |
| `internal/server` | Dependency injection, HTTP server lifecycle |

## Development Notes

- Handlers should only contain HTTP logic.
- Services contain business logic.
- Repositories contain database access.
- Database connection management is isolated in the `database` package.
- Repository interfaces allow services and handlers to be unit tested without a real database.
- Integration tests should always use `TEST_DATABASE_URL`, never the production database.

## Future Improvements

- Authentication / Authorization
- Structured logging (`log/slog`)
- Request middleware
- Health checks (`/healthz`, `/readyz`)
- Prometheus metrics
- OpenTelemetry tracing
- Docker support
- CI/CD pipeline


## Postgres Things

Log in as root user

```bash
sudo -u postgres psql

list
\l


```

Log into a specific database

```bash
psql -U trading_user -d trading_db -h localhost
```