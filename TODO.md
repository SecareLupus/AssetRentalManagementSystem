# Rental Management System - Phase II Roadmap

Building upon the solid backend foundation, the next phases focus on security, user interaction, and operational intelligence.

## Phase 1: Identity & Access Management (IAM)

Secure the API and define user roles to support the varied workflows (Management vs. Technical Operations).

- [x] Define `User` domain model with support for roles (Admin, FleetManager, Technician).
- [x] Implement `UserRepository` and SQL migrations for user tables.
- [x] Implement Authentication Service (JWT-based or OIDC integration).
- [x] Create Authorization Middleware (RBAC) for API endpoints.
- [x] Add `CreatedBy` / `UpdatedBy` audit trails to existing core entities (`RentAction`, `Asset`, etc.).

## Phase 2: Reference UI & API Visualization

Establish the frontend not just as a tool, but as a live reference implementation for the API.

- [x] **OpenAPI Integration**: specific annotations to Go handlers and generate `swagger.json` (using swaggo/swag).
- [x] **Frontend Foundation**: Initialize React/Vite with a "Developer Mode" context.
- [x] **API Inspector Component**: A global UI overlay that listens to network requests and displays:
    -   The exact HTTP method and URL used.
    -   The Request Body / Headers sent.
    -   The Response received.
    -   Relevant documentation excerpt for that endpoint.
- [x] **Dashboard Implementation**: Build the "Commander's Dashboard" utilizing this new Inspector system.

## Phase 3: Catalog & Reservation (Self-Documenting)

Enable users to browse inventory and request equipment (The "Rent" core loop) with full transparency.

- [x] **Catalog View**: Grid/List view of `ItemTypes` with availability status.
- [x] **Asset Details**: Rich view of individual assets (specs, history, maintenance log).
- [x] **Reservation Wizard**: Multi-step form to create a `RentAction` (Select dates, items, logistics).
- [x] **Approval Workflow UI**: Interface for Managers to review, approve, or reject requests.
- [x] **Integration Points**: Ensure every button/action in these views triggers the API Inspector log.

## Phase 4: Fleet Services & Maintenance Station

Tools for the Technician persona to manage the physical lifecycle of devices.

- [x] **Tech Dashboard**: View of "To-Do" items (Inspections due, Provisioning tasks, Returns to process).
- [x] **Inspection Runner**: UI to render the Dynamic Inspection Forms (from Phase 7) and capture results/photos.
- [x] **Provisioning Interface**: Step-by-step wizard for `ProvisionAction` (setting Test Bits, verifying firmware).
- [x] **Check-in/Check-out Kiosk**: Simplified view for scanning assets in and out of the warehouse.

## Phase 5: Planning & Intelligence Engine

Leverage the data to provide predictive insights and reporting.

- [x] **Availability Heatmap**: Visual calendar view showing projected equipment utilization and gaps.
- [x] **Shortage Alerts**: Proactive notifications when overlapping reservations exceed physical inventory.
- [x] **"What-If" Planning Mode**: Ability to draft a large `RentAction` and see its impact on future fleet health without committing.
- [x] **Maintenance Forecasting**: Predicting when assets will need inspection based on usage cycles captured in `Metadata`.ion.
