# Rental Management System - Development Roadmap

This document outlines the planned development phases for the Rental Management System, based on the specifications.

## Phase 1: Foundation & Persistence

Establish the persistence layer and schema definition for upcoming development.

- [ ] Define initial SQL schema (DDL) for core tables (`item_types`, `assets`, `rent_actions`, `rent_action_items`).
- [ ] Implement `SqlRepository` in `internal/db` satisfying the `Repository` interface.
- [ ] Add unit tests for repository methods using a mock or local DB.
- [ ] Setup database connection logic and configuration in `cmd/server`.

## Phase 2: Catalog & Inventory API

Build out the management interfaces for equipment types and physical assets.

- [ ] Implement CRUD handlers for `ItemTypes`.
- [ ] Implement Asset management API (CRUD, status updates).
- [ ] Add validation for `ItemType` schema requirements.
- [ ] Implement "Catalog" view (browsable list of active items).

## Phase 3: Reservation Lifecycle

Implement the core `RentAction` workflow.

- [ ] Implement `CreateRentAction` with validation for duration and item availability.
- [ ] Implement status transition logic (Draft -> Pending -> Approved/Rejected).
- [ ] Add API endpoints for approving and cancelling reservations.
- [ ] Implement `GetRentAction` with nested items and status history.

## Phase 4: Scheduling & Allocation

Add intelligence for asset allocation and planning.

- [ ] Implement allocation logic to bind specific `Assets` to `RentActionItems`.
- [ ] Implement basic conflict detection for overlapping reservations.
- [ ] Add "Shortfall" reporting API to identify missing assets for future needs.
- [ ] Track asset location and assignment changes.

## Phase 5: Event System & Outbox

Implement automation and external integrations.

- [ ] Design and implement the Trigger system for automated workflows.
- [ ] Implement the Outbox pattern for reliable event processing.
- [ ] Create webhook notification system for status changes.
- [ ] Add maintenance logging triggers when assets are returned.
