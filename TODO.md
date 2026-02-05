# Rental Management System - Milestone IV Roadmap

This milestone focuses on refining the user experience and preparing the system for advanced operational flows.

## Phase 16: UX Refinements & Foundations

Addressing immediate usability issues and establishing a better foundation for the UI.

- [x] **Content Pane Z-Layering**: Fix issues where content elements overlap incorrectly.
- [x] **Menu Pane Interaction**: Fix menu overlapping content and add collapsibility.
- [x] **Item Creation**: Add functionality to create new items (UI Placeholder).
- [x] **Reservation UX**: Move "Create Reservation" to the reservations page.
- [x] **General Polish**: Refactored styles, centralized CSS variables, and created shared components.

## Phase 17: Frontend API Full Integration

Goal: Achieve 100% API coverage in the frontend by hooking up all backend calls to their respective UI components.

- [x] **API Audit**: Map all backend endpoints to frontend usage; identify gaps.
- [x] **Catalog Integration**: Implement real `POST` / `PUT` / `DELETE` for Item Types.
- [x] **Inventory Integration**: Implement specific Asset creation, editing, and status management.
- [x] **Reservation Lifecycle**: Implement full flow (Submit -> Approve/Reject -> Provision -> Return).
- [x] **Error Handling**: Standardize API error handling and loading states across the app.

## Future Plans (Carried over from Milestone III)

- [ ] **Remote Management Interface**: Implement a pluggable interface for generic service integration. Connect frontend to Real-time Device controls.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.
