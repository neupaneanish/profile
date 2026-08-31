# Profile Microservice

###### Design and Developed by [Anish Neupane](https://neupaneanish.com.np)

--- 

## Overview

Distributed Profile Microservice with Go, gRPC, PostgreSQL, and Valkey.

---

## Features

- gRPC APIs
- Protocol Buffers (buf.build)
- SQLc query
- PostgreSQL Database
- Valkey for caching
- OpenTelemetry observability
- Dockerized testing (testcontainers)
- Benchmarks, E2E

---

## Technologies Stack

| Technology                                                |                                                                                                  | Description                                                                      |
|:----------------------------------------------------------|:------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------|
| [**Go**](https://go.dev)                                  |             <img src="https://thesvg.org/icons/go/default.svg" height="12" alt="Go">             | Core application logic                                                           |
| [**gRPC**](https://grpc.io)                               |           <img src="https://thesvg.org/icons/grpc/default.svg" height="24" alt="gRPC">           | High-performance RPC framework                                                   |
| [**PostgreSQL**](https://postgresql.org)                  |     <img src="https://thesvg.org/icons/postgresql/default.svg" height="24" alt="PostgreSQL">     | Primary relational database                                                      |
| [**Valkey**](https://valkey.io)                           |         <img src="https://thesvg.org/icons/valkey/default.svg" height="24" alt="Valkey">         | High-performance data structure store                                            |
| [**Docker**](https://docker.com)                          |         <img src="https://thesvg.org/icons/docker/default.svg" height="24" alt="Docker">         | Containerization and deployment                                                  |
| [**Test Containers**](https://testcontainers.com)         | <img src="https://thesvg.org/icons/development-containers/default.svg" height="24" alt="Docker"> | Orchestrates real PostgreSQL and Valkey Docker instances inside automated tests. |
| [**GitHub Actions**](https://github.com/features/actions) | <img src="https://thesvg.org/icons/github-actions/default.svg" height="24" alt="GitHub Actions"> | CI/CD automation pipelines                                                       |
| [**OpenTelemetry**](https://opentelemetry.io)             |  <img src="https://thesvg.org/icons/opentelemetry/default.svg" height="24" alt="OpenTelemetry">  | Observability and telemetry framework                                            |

---

## Endpoints

### External

### Gateway

### Root

---

## Environments

|       Name        |            Default            |            Options            |
|:-----------------:|:-----------------------------:|:-----------------------------:|
|   DATABASE_HOST   |                               |                               |
|   DATABASE_NAME   |                               |                               |
|   DATABASE_USER   |                               |                               |
| DATABASE_PASSWORD |                               |                               |
|   DATABASE_PORT   |            `5432`             |                               |
|   DATABASE_SSL    |            `True`             |                               |
|    VALKEY_URL     |                               |                               |
|       PORT        |            `50051`            |        `80` to `65535`        |
|   SERVICE_NAME    | `neupaneanish.com.np/profile` |                               |
|    ENVIRONMENT    |         `development`         | `development` or `production` |
|   TELEMETRY_URL   |                               |        gRPC port only         |

```dotenv
DATABASE_HOST=
DATABASE_NAME=
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_PORT=
DATABASE_SSL=
VALKEY_URL=127.0.0.1:6379
PORT=50051
SERVICE_NAME=neupaneanish.com.np/profile
ENVIRONMENT=development
TELEMETRY_URL=127.0.0.1:4317
```

---

## Setup, Execution & Testing

```bash
# 1. Clone the core framework engine
git clone https://github.com/neupaneanish/profile.git
cd profile

# 2. Initialize Git submodules
# (Note: if HTTP use git config --global url."https://github.com/".insteadOf "git@github.com:")
git submodule update --init

# 3. Generate Go code from protobuf definitions (Requires Buf CLI)
buf generate

# 4. Generate Go code from SQL queries using SQLc (Requires SQLc CLI)
sqlc generate

# 5. Execute the tests
go test -v -tags=unit ./...
go test -v -tags=integration ./...
go test -v -tags=benchmark ./...
go test -v -tags=e2e ./...

# 6. Launch the local microservice API server
# (Note: Requires an active OpenTelemetry collector instance, e.g., SigNoz)
go run cmd/server/main.go
```

---

## Coverage ~96.00%

> Note: Metrics reflect core application logic after filtering out `main.go`, generated protobuf definitions, raw SQL
> repository code, and test helper suites.

> Coverage is done through real infrastructure PostgreSQL, Valkey, OpenTelemetry i.e. testcontainers. It doesn't have
> any mocks.

```bash
# Generate coverage
go test -race -tags=unit,integration,benchmark,e2e -coverprofile=coverage.out -coverpkg=./... ./..

# Filter out external boundaries, generated code, and tooling 
grep -v -E "cmd/|/internal/protobuf/|/internal/repository/|/tests/|/protobuf/|/database/" coverage.out > coverage.clean.out

# Export to interactive HTML for local branch analysis
go tool cover -html=coverage.clean.out -o coverage.clean.html 

# Output statement breakdown to CLI
go tool cover -func=coverage.clean.out 
```

---

## Testing Architecture (Testcontainers)

This repository uses a modern, completely containerized testing environment:

- **Integration and E2E Tests:** Used real database, valkey and telemetry instances for integration tests.
- **Benchmark Tests:** Used memory server i.e. `bufconn` instead of real server for tests.

---

## [License](LICENSE)