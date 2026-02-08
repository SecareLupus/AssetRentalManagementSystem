# Rental Management System - Milestone VII Roadmap

This milestone focuses on cleaning up temporary assumptions, fleshing out our interaction models, and cleaning up the UI for presentation.

## Phase 27: System Stability & Stats Alignment
Goal: Resolve critical 500 errors and align dashboard statistics with modernized logistics models.

- [x] **Fix Places API**: Handle NULL `presumed_demands` in `sql_repository.go` to resolve 500 errors.
- [x] **Modernize Dashboard Stats**: Update `GetDashboardStats` to use `rental_reservations` and `check_out_actions` instead of legacy `rent_actions`.
- [x] **Scorecard Grid**: Refactor Dashboard UI to a consistent 2x3 grid layout.
- [x] **Fleet Reports Routing**: Fix routing to prevent blank pages when accessing reports.

## Phase 28: Asset Identity & Lifecycle
Goal: Refine how assets are identified, tracked, and initialized.

- [x] **Identifying Codes Review**: Audit SKUs, serial numbers, and product codes to prevent conflation.
- [x] **Component Tracking**: Implement logic for tracking internal component serial numbers during refurbishment.
- [x] **Default Internal Location**: Establish a "Default Internal Location" and assign it to any asset created without a location.
- [x] **Conditional UI Fields**: Hide/show fields (like serial numbers) in the UI based on Item Type `supported_features`.

## Phase 29: Entity Management & Optimization
Goal: Fix broken UI links and improve data entry workflows for personnel and templates.

- [x] **Location Dropdown Fix**: Ensure the Asset creation dropdown correctly populates from the Places API.
- [x] **Personnel Edit Fix**: Resolve issues with the Personnel edit page failing to load/save email, phone, and role.
- [x] **View Profile Implementation**: Flesh out the "View Profile" link for Personnel.
- [x] **Inspection Template Polish**: Fix the template edit page and implement a way to assign templates to Item Types.

## Phase 30: System Admin & IAM
Goal: Provide UI-driven management for users, roles, and global settings.

- [ ] **User Settings Page**: Implement a dedicated page for individual user preferences.
- [ ] **Server Admin Page**: Create an administrative interface for IAM (User/Directory management).
- [ ] **Configurable Place Types**: Allow system-level configuration of place types and their data subsets.

## Future Plans

- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.
- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **Third-Party Logistics (3PL) Integration**: Connect with shipping partners for automated dispatch.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.
