# Rental Management System - Milestone VI Roadmap

This milestone focuses on closing the loop between deployment, to target, to maintenance, to inspection, to re-deployment.

## Phase 23: Entity Management

Goal: Implement first class entities to represent the deployment targets: Companies, Company Contacts, Company Sites (Facilities), Company Locations (Rooms), Company Events (Projects), Event Asset Needs.

### Companies

- [x] **Company Management**: Implement a UI for creating and managing companies.
- [x] **Company Contact Management**: Implement a UI for creating and managing company contacts. Contacts represent individuals who may be associated with a company.
- [x] **Company Site Management**: Implement a UI for creating and managing company site. Sites represent facilities utilized by a company which possess a mailing address.
- [x] **Company Location Management**: Implement a UI for creating and managing company locations. Locations represent physical spaces within a company site. Locations may be nested to represent subspaces. Locations may have presumed asset needs, which will be applied to an Event Asset Need Calculation automatically (unless overridden), if that location is included in the event.
- [x] **Company Event Management**: Implement a UI for creating and managing company events. Events have a start and end date, belong to a company, and may reference any number of Contacts, Sites, Locations, and Asset Needs which are disconnected from a given location.
- [x] **Company Event Asset Need Management**: Implement a UI for creating and managing company asset needs for each event.
- [x] **Company, Site, Location, Event, and Asset Need API Integration**: Implement a pluggable interface for integration with generic outside REST services.

## Phase 24: Structural Consolidation (Places & Roles)

Goal: Converge the physical and organizational data models toward Schema.org standards to support infinite nesting and multi-faceted relationships.

### Unified Place Model
- [x] **Recursive Place Schema**: Consolidate `Site` and `Location` into a single recursive `Place` entity. Add `contained_in_place_id` and optional `PostalAddress` fields.
- [x] **Data Migration (Sites to Places)**: Execute SQL migrations to port facility data into the unified model while maintaining referential integrity for existing assets.
- [x] **Hierarchical Entity Manager**: Update the `EntityManager` UI to support an arbitrary depth of nested places (Site > Building > Room > Cabinet).

### Personnel & Organizations
- [x] **Person/Organization Decoupling**: Refactor `Contact` as a standalone `Person` entity that can have multiple `OrganizationRole` relationships.
- [x] **ContactPoint Implementation**: Implement `ContactPoint` to manage communication lines (email/phone) scoped to specific organizations or events.


## Phase 25: Logistics Modernization (Reservations & Demands)

Goal: Transition from monolithic action records to a more flexible Reservation/Demand ecosystem for predictive logistics and partial fulfillment.

### Intent & Fulfillment
- [x] **RentalReservation Model**: Implement a first-class `RentalReservation` entity to track "intent to rent" separate from the act of movement.
- [x] **Granular Action Logging**: Implement specific `CheckOutAction` and `ReturnAction` entities linked to a parent Reservation.
- [ ] **Partial Fulfillment Engine**: Update fulfillment logic to allow multiple check-out events for a single reservation (staggered delivery).

### Offer/Demand Ecosystem
- [x] **Logistics Demand Model**: Replace `EventAssetNeed` with the `Demand` model, incorporating `businessFunction` and `eligibleDuration`.
- [x] **Standardized Vocabulary**: Align API JSON tags with Schema.org vocabulary (camelCase, JSON-LD context) while providing a compatibility layer for internal legacy IDs.
- [x] **Pre-Fulfillment Validation**: Implement standard validation logic that checks `Asset` availability against global `Demand` timelines across all facilities.

## Phase 26: Fulfillment Orchestration

Goal: Implement the "Engine" that manages the delta between intent (Reservations) and reality (Movements).

### Workflow & Logic
- [ ] **Fulfillment State Machine**: Implement service logic to calculate "Remaining Needs" by comparing Demands vs. CheckOutActions.
- [ ] **Staggered Checkout Workflow**: Implement API to batch-dispatch assets, auto-generating `CheckOutAction` records and updating `Asset` statuses in one transaction.
- [ ] **Status Auto-Promotion**: Implement logic to transition `RentalReservation` status to `PartiallyFulfilled` or `Fulfilled` based on actual asset movement.
- [ ] **Return Pipeline**: Implement the return-to-warehouse workflow that reconciles `CheckOutActions` with `ReturnActions`.

## Future Plans

- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.
- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **Third-Party Logistics (3PL) Integration**: Connect with shipping partners for automated dispatch.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.
