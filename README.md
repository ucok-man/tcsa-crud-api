# Technical Test PT Cipta Solusi Aplikasi | Transaction Management API

A RESTful API for managing financial transactions with filtering, sorting, and analytics capabilities built with Go, Echo framework, and PostgreSQL.

## Features

- **CRUD Operations**: Create, read, update, and delete transactions
- **Transaction Status Management**: Handle pending, success, and failed states
- **Filtering & Pagination**: Filter by status, user ID, date range with customizable page sizes
- **Analytics Dashboard**: Transaction summary with status distribution and rate percentages
- **OpenAPI Documentation**: Interactive Swagger UI for API exploration
- **Optimistic Locking**: Prevent concurrent modification conflicts

## Tech Stack

- **Language**: Go 1.25.0
- **Framework**: Echo v4
- **Database**: PostgreSQL 17
- **Migration**: Goose
- **Validation**: go-playground/validator
- **Logging**: Uber Zap

## Prerequisites

- Go 1.25.0 or higher
- PostgreSQL 17
- Docker & Docker Compose (optional)
- Make (Linux/macOS users, for convenience commands)
- Goose (For running migration)
- Air (For hot reload)

## Quick Start

### Installing Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

# This will install the goose binary to your $GOPATH/bin directory.
```

### Installing Air

```bash
# With go 1.25 or higher:
go install github.com/air-verse/air@latest

# This will install the air binary to your $GOPATH/bin directory.
```

### Option 1: Using Docker

1. **Pull the Docker image**

```bash
docker pull ucokman/tsca_api:latest
```

2. **Start PostgreSQL database**

```bash
docker run -d \
  --name tcsa_postgres \
  -e POSTGRES_DB=tcsa \
  -e POSTGRES_USER=tcsa \
  -e POSTGRES_PASSWORD=pa55word \
  -p 5433:5432 \
  postgres:17
```

3. **Run migrations**

```bash
# Run migrations
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="postgres://tcsa:pa55word@localhost:5433/tcsa?sslmode=disable" \
goose --dir ./migrations up

# Or via makefile
make db/up
```

4. **Start the API container**

```bash
docker run -d --name tcsa_api \
  -p 4000:4000 \
  -e TCSA_DB_DSN="postgres://tcsa:pa55word@tcsa_postgres:5433/tcsa?sslmode=disable" \
  ucokman/tsca_api:latest
```

The API will be available at `http://localhost:4000`

5. **To see list all available options**

```bash
go run ./cmd/api --help

# Or
make api/help
```

### Option 2: Using Docker Compose (Recommended)

1. **Clone the repository**

```bash
git clone git@github.com:ucok-man/tcsa.git
cd tcsa
```

2. **Start all services with migrations**

```bash
make compose/up
# When prompted, type 'y' to run migrations first
```

The API will be available at `http://localhost:4000`

3. **To clean up all data**

```bash
make compose/clear
```

### Option 3: Running Manual

1. **Clone the repository**

```bash
git clone git@github.com:ucok-man/tcsa.git
cd tcsa
```

2. **Install dependencies**

```bash
go mod download
go mod vendor
```

3. **Set up environment variables**

```bash
cp .env.example .env.development
# Edit .env.development with your database credentials
```

4. **Run migrations**

```bash
# Run migrations
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="postgres://tcsa:pa55word@localhost:5433/tcsa?sslmode=disable" \
goose --dir ./migrations up
```

5. **Start the application**

For development:

```bash
air -c .air.toml # hot reload if installed

# Or
go run ./cmd/api/...
```

Or build and run:

```bash
go build -ldflags='-s' -o=./bin/api ./cmd/api/ && ./bin/api
```

The API will be available at `http://localhost:4000`

## Database Migrations

### Running Migrations (with Make)

```bash
# Apply all pending migrations
make migrate/up

make migrate/down
# Confirm with 'y' when prompted
```

### Manual Migration (without Make)

```bash
# Apply migrations
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="postgres://tcsa:pa55word@localhost:5432/tcsa?sslmode=disable" \
goose --dir ./migrations up

# Rollback all migrations
GOOSE_DRIVER=postgres \
GOOSE_DBSTRING="postgres://tcsa:pa55word@localhost:5432/tcsa?sslmode=disable" \
goose --dir ./migrations reset
```

## API Documentation

Once the application is running, visit:

- **Swagger UI**: http://localhost:4000/docs
- **OpenAPI Spec**: http://localhost:4000/swagger.yaml
- **Health Check**: http://localhost:4000/healthcheck

## Available Endpoints

### Health

- `GET /healthcheck` - API health status

### Transactions

- `GET /transactions` - Get all transactions (with pagination, filtering, sorting)
- `POST /transactions` - Create a new transaction
- `GET /transactions/:id` - Get transaction by ID
- `PUT /transactions/:id` - Update transaction
- `DELETE /transactions/:id` - Delete transaction

### Dashboard

- `GET /dashboard/summary` - Get transaction summary and analytics

## Environment Variables

| Variable                    | Description                                       | Default            |
| --------------------------- | ------------------------------------------------- | ------------------ |
| `TCSA_PORT`                 | Server port                                       | `4000`             |
| `TCSA_ENV`                  | Environment (development/staging/production)      | `development`      |
| `TCSA_DB_DSN`               | PostgreSQL connection string                      | See `.env.example` |
| `TCSA_DB_MAX_OPEN_CONN`     | Maximum open database connections                 | `25`               |
| `TCSA_DB_MAX_IDLE_CONN`     | Maximum idle database connections                 | `15`               |
| `TCSA_DB_MAX_IDLE_TIME`     | Maximum idle time for connections (time.Duration) | `15m`              |
| `TCSA_LOG_LEVEL`            | Logging level (debug/info/warn/error)             | `debug`            |
| `TCSA_CORS_TRUSTED_ORIGINS` | Allowed CORS origins (comma-separated)            | `""`               |

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with gotestdox (if installed)
make test/doc

# Run tests with coverage
go test -v -cover ./...
```

### Code Quality

```bash
# Format code, tidy dependencies, and vendor
make tidy

# Run quality checks (vet, staticcheck, tests)
make audit
```

## Database Management

```bash
# Start database only
make db/up

# Stop database
make db/down

# Remove database with volumes
make db/clear

# Wait for database to be ready
make db/wait
```

## Project Structure

```
.
├── cmd/api/             # Application entry point
│   ├── main.go          # Main application file
│   ├── config.go        # Configuration management
│   ├── routes.go        # Route definitions
│   ├── handler_*.go     # HTTP handlers
│   ├── middleware.go    # Custom middleware
│   └── docs/            # Swagger documentation
├── internal/
│   ├── data/            # Data models and database logic
│   ├── validator/       # Request validation
│   ├── serializer/      # JSON serialization
│   ├── tlog/            # Logging wrapper
│   └── utility/         # Helper functions
├── migrations/          # Database migrations
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose configuration
├── Makefile             # Build and deployment commands
└── go.mod               # Go module dependencies
```
