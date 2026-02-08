# Rental Management System - Milestone VIII Roadmap

This milestone focuses on integration with outside sources of truth, enabling the system to sync with enterprise tools, identity providers, and logistics partners.

## Phase 31: External Identity & User Preferences
Goal: Modernize user management with self-service preferences and support for corporate SSO.

- [ ] **Individual User Preferences**: Implement a dedicated page for personal user settings (Timezone, Notifications, Profile).
- [ ] **OIDC/SAML Provider Integration**: Research and implement support for external identity providers (e.g., Okta, Google Workspace).

## Phase 32: Enterprise ERP & Inventory Sync
Goal: Connect internal inventory with external enterprise resource planning systems.

- [ ] **ERP Integration Layer**: Create a generic connector for syncing assets and item types with outside ERP systems.
- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.

## Phase 33: Logistics & 3PL Integration
Goal: Automate shipping and dispatch with third-party logistics partners.

- [ ] **3PL Dispatch automation**: Connect with shipping partners (UPS/FedEx/DHL) for automated label generation and dispatch.
- [ ] **Real-time Transit Tracking**: Sync shipment statuses into the Reservation / Dispatch UI.

## Future Plans

- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.