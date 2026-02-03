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
- SQL Database (PostgreSQL recommended)

### Installation

```bash
go mod download
```

### Running the Service

```bash
go run ./cmd/server
```

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
