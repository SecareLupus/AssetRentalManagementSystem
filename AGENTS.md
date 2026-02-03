# Agent Guidelines: Rental Management System

This document provides context and guidelines for AI agents (like Antigravity) working on the Rental Management System.

## Project Context

The Rental Management System is an API-first service designed to manage rental reservations (RentActions), catalog item types, assets, and stock pools. It aligns with Schema.org standards (specifically `RentAction`).

## Technical Stack

- **Language**: Go
- **Architecture**: Clean Architecture / Standard Go project layout
- **Database**: SQL (PostgreSQL preferred)
- **API**: JSON-LD / Schema.org aligned

## Code Quality & Standards

- Follow standard Go idioms and naming conventions.
- Use the standard Go project layout:
  - `cmd/`: Application entry points.
  - `internal/`: Private code not intended for public use.
  - `pkg/`: Public library code.
  - `api/`: API definitions.
- Ensure all `RentAction` payloads are Schema.org compliant where possible.
- Use the Outbox pattern for all trigger/event emissions.

## Working with the Schema

Refer to the `Rental Reservations System Spec.pdf` for the core table definitions. When modifying the database schema, ensure migrations are created in the `migrations/` directory.

## Triggers & Webhooks

The system uses a trigger system for side effects. Ensure new event types are documented in the specification and followed in the implementation.
