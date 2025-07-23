# GoFr Account Service

A RESTful user account management service built with GoFr framework and PostgreSQL.

## Quick Start

### Using Docker (Recommended)

```bash
# Clone and navigate to the project

# Start the application with database
docker-compose up --build

# The API will be available at http://localhost:8080
```

### Local Development

```bash
# Setup PostgreSQL
CREATE USER gouser WITH PASSWORD 'pass';
CREATE DATABASE accounts;
GRANT ALL PRIVILEGES ON DATABASE accounts TO gouser;

# Apply migration
psql -U gouser -d accounts -f internal/db/migrations/001_create_users.sql

# Run the application
go run cmd/main.go
```

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### User Registration
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"gautam","email":"gautam@email.com","password":"pass"}'
```

### Get User
```bash
# By username
curl "http://localhost:8080/user?username=gautam"

# By email
curl "http://localhost:8080/user?email=gautam@email.com"

# By ID
curl "http://localhost:8080/user?id=1"
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./test -run TestRegisterUser
```

## Features

- User registration with validation
- Secure password hashing (bcrypt)
- User retrieval by username/email/ID
- PostgreSQL database integration
- Docker containerization
- Comprehensive unit tests
- CI/CD pipeline (GitHub Actions)
- RESTful API design

## Tech Stack

- **Framework**: GoFr
- **Language**: Go 1.23
- **Database**: PostgreSQL 14
- **Authentication**: bcrypt password hashing
- **Testing**: Go testing framework with mocks
- **Containerization**: Docker & Docker Compose

## Project Structure

```
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── handler/               # HTTP handlers
│   ├── service/               # Business logic
│   ├── repository/            # Data access layer
│   ├── model/                 # Data models
│   └── db/migrations/         # Database migrations
├── pkg/hashutil/              # Utility packages
├── test/                      # Unit tests
└── docker-compose.yml         # Container orchestration
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_USER | Database user | gouser |
| DB_PASSWORD | Database password | pass |
| DB_NAME | Database name | accounts |
| APP_PORT | Application port | 8080 |

## Deliverables

This project includes:

- Project repository with Go modules, Docker-Compose, and CI workflow
- Database schema and migrations for users table
- User registration and retrieval endpoints
- Secure password hashing using bcrypt
- Unit tests for all endpoints
- Sample requests and documentation
