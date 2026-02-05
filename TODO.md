# Rental Management System - Milestone III Roadmap

Building on the foundation established in Milestone II, this phase focuses on hardening system integrity, expanding real-time connectivity, and enabling advanced operational flows.

## Phase 12: M-II Hardening & Security

Address critical integration gaps identified during the Milestone II audit to ensure system production-readiness.

- [x] **API Security**: Apply `AuthMiddleware` to all `/v1` routes in `router.go`.
- [x] **Event Instrumentation**: Update core handlers (`CreateAsset`, `UpdateAsset`, `ApproveRentAction`, etc.) to call `AppendOutboxEvent`.
- [x] **Audit Trail Expansion**: Add `CreatedByUserID` and `UpdatedByUserID` to the `ItemType` model and ensure `UpdatedByUserID` is captured across all core entities.
- [x] **Webhook Dispatcher**: Implement a generic HTTP dispatcher within the `OutboxWorker` to deliver events to registered `WebhookConfigs`.

## Phase 13: MQTT Integration & Real-time Mirroring

Introduce MQTT as a primary event conduit for edge clients and mobile observers.

- [x] **MQTT Infrastructure**: Add an MQTT Client to the backend (e.g., using `paho.mqtt.golang`).
- [x] **Outbox-to-MQTT Mirror**: Implement an MQTT adapter in `OutboxWorker` that echoes every processed event to a structured topic tree (e.g., `rms/events/{event_type}`).
- [x] **Health Status Mirroring**: Periodically publish asset health summaries obtained via `RemoteManager` to MQTT.

## Phase 14: Fleet Connectivity & Remote Ops

Extend the ability to interact with and verify remote hardware.

- [ ] **RemoteManager Implementation**: Implement the first concrete provider (e.g., a MeshCentral or SSH-based agent).
- [ ] **Real-time Dashboard Metrics**: Integrate MQTT or long-polling into the Dashboard to show live device health without page refreshes.
- [ ] **Power Action Verification**: Ensure `ApplyPowerAction` results are captured as events and reflected in the asset history.

## Phase 15: Advanced Lifecycle & Logistics

Optimizing the flow of assets through the facility and field.

- [ ] **Asset Reclamation UI**: Dedicated interface for bulk-recalling assets (e.g., "End of Project" wizard).
- [ ] **Inventory Reconciliation**: A "Scan & Compare" tool for the Warehouse Kiosk to verify database records against physical inventory.
- [ ] **Maintenance Prediction refinement**: Tuning the Intelligence Engine with more granular usage metrics.

## Future Plans (Backlog)

- **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.
