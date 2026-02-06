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

## Phase 18: Full System Readiness

Goal: Bridge the gap from "demo-capable" to "operation-ready" by closing API gaps, removing placeholders, and refining operational utility.

- [x] **100% API Coverage & Specialized UI**
  - [x] **Admin Operations Hub**: Implement Bulk Recall and Inventory Reconciliation views.
  - [x] **Lifecycle Management**: Implement Archival/Soft-Delete for Item Types in the Catalog.

- [x] **Data Foundations (Eliminate Placeholders)**
  - [x] **Auth Integration**: Link all submission metadata directly to the authenticated Session/User ID.
  - [x] **Reactive Status Tracking**: Replace static timers with real-time polling of Asset `provisioning_status`.
  - [x] **Live System Health**: Connect dashboard health widgets to real system and service telemetry.

- [x] **Operational Use Case Completion**
  - [x] **"What-If" Planning Simulator**: Enable interactive capacity testing and shortage simulation in the Intelligence Hub.
  - [x] **High-Efficiency Batching**: Optimize `WarehouseKiosk` for high-speed barcode scanning with hands-free logic.
  - [x] **Fleet Utilization Metrics**: Implement advanced reporting for downtime trends and refurbishment efficiency.
