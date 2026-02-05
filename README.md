# Rental Management System

A standalone Reservations Service supporting planning and allocation of fleet equipment.

## Overview

This service is API-first, providing robust management for rental reservations (RentActions), cataloging item types, and tracking physical assets. It integrates with Schema.org standards to ensure interoperability.

### Key Features

- **Schema.org Alignment**: Uses `RentAction` for reservations and other standard types where applicable.
- **Planning-Aware**: Predicts future device needs and surfaces shortfalls.
- **Trigger System**: Configurable automation via webhooks and internal actions using the Outbox pattern.
- **Maintenance Logging**: Comprehensive logs for asset maintenance history.

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (for Postgres/MQTT)

### Setup Test Environment

The fastest way to spin up the entire stack (App, Database, and MQTT broker) is via Docker Compose:

```bash
docker-compose up -d --build
```

This will:

1. Build the Go application container.
2. Start a PostgreSQL database instance.
3. Start a Mosquitto MQTT broker.
4. Expose the API at `http://localhost:8080`.

### Configuration

The following environment variables can be used to configure the service. You can create a `.env` file based on `.env.example` for local overrides:

```bash
cp .env.example .env
```

| Variable            | Description                          | Default                                                                 |
| ------------------- | ------------------------------------ | ----------------------------------------------------------------------- |
| `DATABASE_URL`      | PostgreSQL connection string         | `postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable` |
| `MQTT_BROKER`       | MQTT broker URL                      | `tcp://localhost:1883`                                                  |
| `POSTGRES_USER`     | DB User (used by docker-compose)     | `postgres`                                                              |
| `POSTGRES_PASSWORD` | DB Password (used by docker-compose) | `postgres`                                                              |
| `POSTGRES_DB`       | DB Name (used by docker-compose)     | `rental_db`                                                             |

### Database Migrations

Apply the SQL migrations located in `./migrations/` to your database. You can do this manually using `psql`:

```bash
for f in migrations/*.sql; do psql "$DATABASE_URL" -f "$f"; done
```

### Running the Service

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Run the server:

   ```bash
   go run ./cmd/server
   ```

3. Explore the API via Swagger UI:
   Navigate to [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

## Project Structure

- `api/`: API definitions and documentation.
- `cmd/`: Application entry points (server, workers).
- `internal/`: Core business logic and database access.
- `pkg/`: Shared libraries and utilities.
- `migrations/`: SQL migration files.
- `scripts/`: Development and deployment scripts.

## Documentation

Refer to [Rental Reservations System Spec.pdf](./Rental%20Reservations%20System%20Spec.pdf) for the full technical specification.
For AI agents, see [AGENTS.md](./AGENTS.md).
