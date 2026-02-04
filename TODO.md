# Rental Management System - Development Roadmap

This document outlines the planned development phases for the Rental Management System, based on the specifications.

## Phase 1: Foundation & Persistence

Establish the persistence layer and schema definition for upcoming development.

- [x] Define initial SQL schema (DDL) for core tables (`item_types`, `assets`, `rent_actions`, `rent_action_items`) with Schema.org alignment.
- [x] Implement `SqlRepository` in `internal/db` satisfying the `Repository` interface.
- [x] Add unit tests for repository methods using a mock or local DB.
- [x] Setup database connection logic and configuration in `cmd/server`.

## Phase 2: Catalog & Inventory API

Build out the management interfaces for equipment types and physical assets.

- [x] Implement CRUD handlers for `ItemTypes`.
- [x] Implement Asset management API (CRUD, status updates).
- [x] Add validation for `ItemType` schema requirements.
- [x] Implement "Catalog" view (browsable list of active items).

## Phase 3: Reservation & Deployment Lifecycle

Implement the core `RentAction` workflow, mapping it to the "Deploy" phase of the fleet lifecycle.

- [x] Implement Submit, Approve, Reject, and Cancel logic in `RentAction` domain model.
- [x] Add status transition handlers to the API.
- [x] Implement basic availability checks during the approval process.
- [x] Ensure appropriate timestamps (`ApprovedAt`, etc.) are captured.
- [x] Verify state transitions with unit tests.

## Phase 4: Fleet Provisioning & Build Specs

Handle the "Provision" and "Inspect" steps of the lifecycle.

- [x] Add fleet-specific fields to `Asset` (BuildSpec, Firmware, ProvisioningStatus, Hostname).
- [x] Implement `BuildSpec` management (defining hardware/software standards).
- [x] Create `ProvisionAction` API for tracking device preparation.
- [x] Implement Build Spec compliance tracking (Inspect -> Verify test bits).

## Phase 5: Refurbishment & Maintenance [REVISED]

Manage the return, repair, and upgrade circular workflow.

- [x] Implement "Recall" workflow (bulk transitions from `deployed` to `recalled`).
- [x] Add `MaintenanceLog` and `RepairHistory` tracking for assets.
- [x] Implement "Refurbish" workflow (Upgrading recalled devices to latest Build Spec).
- [x] Track "Test Bits" and final QC approval before re-entering the pool.

## Phase 6: Event System & External Integrations

- [ ] Generic Remote Management integration (abstracted from specific providers like MeshCentral).
- [ ] Outbox pattern for syncing with SnipeIT and InvenTree.
- [ ] Webhook/Trigger system for automated lifecycle transitions.
- [ ] TODO: Integrate RemoteManagementID into business logic (e.g. auto-recall on health failure).

## Phase 7: Dynamic Inspection Forms

- [x] Define `InspectionTemplate` schema with support for Boolean, Text, and Image fields.
- [x] Implement assignment of templates to `ItemTypes`.
- [x] Create API for collating necessary inspections for a device based on its type.
- [x] Implement submission and storage of inspection results.
