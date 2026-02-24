# Rental Management System - Milestones

## Milestone VIII: Logistics Workflow & Season Planning (Porting Prototype Features)

This milestone focuses on integrating the highly contextualized business workflows validated in the `SGL-Equipment-Logistics` prototype into our generalized Go/React platform. The goal is to keep the Company Backend as harvested via the Universal Ingestion Engine as the source of truth while adding the "Season to Shipment" business logic.

### Phase 31: Season, Show, and Ring Hierarchy
Goal: Implement the hierarchy of `ShowCompany` -> `Season` -> `Show` -> `Ring` and map equipment load-outs per ring on top of our generic models.
- [x] **Backend Schema Overlay**: Add local database tables/structs for `ShowCompany`, `Season`, `Show`, and `Ring` that reference our ingested Source of Truth data.
- [x] **Equipment Load-outs**: Define standard Asset load-outs per Ring to logically generate `Demand` and `RentalReservation` requirements.
- [x] **Predictive Planning API**: Build Go endpoints that predict next year's equipment requirements for a Show based on the previous year's historical data, supporting manual overrides.
- [x] **Frontend - Show & Season Planner**: Build a visual dashboard to create and track the Season hierarchy, visualize predicted Ring loads, and confirm or override equipment demands.

### Phase 32: Delivery & Shipment Routing
Goal: Port the explicit `ScheduledDelivery` and `Shipment` tracking from the prototype into our `CheckOutAction`/`ReturnAction` flow.
- [x] **Delivery Models**: Create `ScheduledDelivery` and `Shipment` structs in Go to group `CheckOutActions` logically.
- [x] **Carrier Tracking Integration**: Add tracking number and carrier fields to the new Shipment models.
- [x] **Frontend - Logistics Dispatch Board**: Create a React UI dedicated to assembling Shipments from approved Reservations and dispatching them.
- [x] **Frontend - Return Processing**: Create a streamlined UI for receiving `ReturnShipments` and instantly triggering the `InspectionRunner`.

### Phase 33: Contextual Device Allocation
Goal: Simplify the UX by allowing users to allocate specific devices (Assets) to Shipments based on the prototype's leaner workflow.
- [ ] **Asset Allocation API**: Create streamlined Go endpoints that assign an `Asset` to a `Shipment` (which under the hood creates the complex `CheckOutAction` and `Demand` lines).
- [ ] **Frontend - Allocation UI**: Build a drag-and-drop or checklist React component for assigning available Assets to pending Shipments.

---

## Milestone IX: External Integrations

This milestone focuses on integration with outside sources of truth, enabling the system to sync with enterprise tools, identity providers, and logistics partners.

### Phase 34: External Identity & User Preferences
Goal: Modernize user management with self-service preferences and support for corporate SSO.
- [ ] **Individual User Preferences**: Implement a dedicated page for personal user settings (Timezone, Notifications, Profile).
- [ ] **OIDC/SAML Provider Integration**: Research and implement support for external identity providers (e.g., Okta, Google Workspace).

### Phase 35: Enterprise ERP & Inventory Sync
Goal: Connect internal inventory with external enterprise resource planning systems.
- [ ] **ERP Integration Layer**: Create a generic connector for syncing assets and item types with outside ERP systems.
- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.

### Phase 36: Logistics & 3PL Integration
Goal: Automate shipping and dispatch with third-party logistics partners.
- [ ] **3PL Dispatch automation**: Connect with shipping partners (UPS/FedEx/DHL) for automated label generation and dispatch.
- [ ] **Real-time Transit Tracking**: Sync shipment statuses into the Reservation / Dispatch UI.

---

## Future Plans

- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.