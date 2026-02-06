# Rental Management System - Milestone V Roadmap

This milestone focuses on fixing issues identified during testing and preparing for the next steps of development.

## Phase 19: Operational Integrity (Critical Fixes)

Goal: Fix broken core logic and API integrations to ensure the system is functionally sound.

- [ ] **Dashboard Interactivity**: Fix non-interactive category list items and "View All" button.
- [ ] **Catalog Core Fixes**: Fix "Filters" and "Add to Cart" buttons.
- [ ] **Item Type Actions**: Fix "Request Reservation" button and make asset list interactive.
- [ ] **Reservation Wizard Logic**: Fix quantitiy accumulation (avoid duplicate rows) and resolve reservation creation failures.
- [ ] **Maintenance Submission**: Resolve failures in submitting inspections.

## Phase 20: Usability & UX Refinements

Goal: Smooth out user workflows and improve data presentation clarity.

- [ ] **Wizard Defaults**: Set default dates to today + 7 days; clarify "ASAP" priority.
- [ ] **Warehouse Inputs**: Add start/end dates and location fields for bulk transactions.
- [ ] **Intelligence Hub Parameters**: Make Critical Shortage and Forecast Frequency editable in Item Type settings.
- [ ] **Heatmap Polish**: Fix pagination buttons and ensure bulk checkouts (without schedules) are respected.
- [ ] **Catalog Archive Cleanup**: Hide the archive button for already archived item types.

## Phase 21: Advanced Management Tooling

Goal: Implement specialized editors and complex simulation logic.

- [ ] **Feature Management**: Enable editing of "Supported Features" and implement functional logic for assigned features.
- [ ] **Inspection Template Editor**: Implement a UI for creating and managing custom inspection templates.
- [ ] **Scannable Tag Editor**: Implement alias-to-asset mapping with regex support for third-party tags (QR extraction).
- [ ] **Simulator Enhancements**: Fix "Launch Simulator" button and update logic to treat assets from non-overlapping scenarios as available.
- [ ] **Forecast Controls**: Implement a "Snooze" option for service forecasts on devices stuck in the field.

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
