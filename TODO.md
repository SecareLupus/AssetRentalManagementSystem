# Rental Management System - Milestone V Roadmap

This milestone focuses on fixing issues identified during testing and preparing for the next steps of development.

## Phase 19: Operational Integrity (Critical Fixes)

Goal: Fix broken core logic and API integrations to ensure the system is functionally sound.

- [x] **Dashboard Interactivity**: Fix non-interactive category list items and "View All" button.
- [x] **Catalog Core Fixes**: Fix "Filters" and "Add to Cart" buttons.
- [x] **Item Type Actions**: Fix "Request Reservation" button and make asset list interactive.
- [x] **Reservation Wizard Logic**: Fix quantitiy accumulation (avoid duplicate rows) and resolve reservation creation failures (SQL `INSERT` column mismatch).
- [x] **Maintenance Submission**: Resolve failures in submitting inspections (Foreign Key violation `fk_is_template`).

## Phase 20: Usability & UX Refinements

- [x] Reservation Wizard: Remove redundant ASAP checkbox (Start date defaults to today)
- [x] Item Type Details: Hide archive button for archived items (Archive/Restore toggle)
- [x] Warehouse Kiosk: Add Destination and Est. Return Date fields for bulk checkouts
- [x] Item Type Settings: Make Critical Shortage and Forecast Horizon editable
- [x] Availability Heatmap: Implement pagination (weekly navigation) and subtract ad-hoc usage from availability logic (without schedules) are respected.
- [x] **Catalog Archive Cleanup**: Hide the archive button for already archived item types.

## Phase 21: Advanced Management Tooling

Goal: Implement specialized editors and complex simulation logic.

- [x] **Feature Management**: Enable editing of "Supported Features" and implement functional logic for assigned features.
- [x] **Inspection Template Editor**: Implement a UI for creating and managing custom inspection templates.
- [x] **Scannable Tag Editor**: Implement alias-to-asset mapping with regex support for third-party tags (QR extraction).
- [x] **Simulator Enhancements**: Fix "Launch Simulator" button navigation and update logic to treat assets from non-overlapping scenarios as available.
- [x] **Forecast Controls**: Implement a "Snooze" option for service forecasts on devices stuck in the field.

## Phase 22: Enterprise Connectivity & Specialized Personas

Goal: Extend the system reach into deep device management and field operations.

- [ ] **Remote Management Interface**: Implement a pluggable interface for generic service integration. Connect frontend to Real-time Device controls.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.

## Future Plans

- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.
- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **Third-Party Logistics (3PL) Integration**: Connect with shipping partners for automated dispatch.
