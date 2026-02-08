# Rental Management System - Milestone VII Roadmap (ARCHIVED)

This milestone focused on cleaning up temporary assumptions, fleshing out our interaction models, and cleaning up the UI for presentation.

## Phase 27: System Stability & Stats Alignment [DONE]
- [x] **Fix Places API**: Handle NULL `presumed_demands` in `sql_repository.go` to resolve 500 errors.
- [x] **Modernize Dashboard Stats**: Update `GetDashboardStats` to use `rental_reservations` and `check_out_actions` instead of legacy `rent_actions`.
- [x] **Scorecard Grid**: Refactor Dashboard UI to a consistent 2x3 grid layout.
- [x] **Fleet Reports Routing**: Fix routing to prevent blank pages when accessing reports.

## Phase 28: Asset Identity & Lifecycle [DONE]
- [x] **Identifying Codes Review**: Audit SKUs, serial numbers, and product codes to prevent conflation.
- [x] **Component Tracking**: Implement logic for tracking internal component serial numbers during refurbishment.
- [x] **Default Internal Location**: Establish a "Default Internal Location" and assign it to any asset created without a location.
- [x] **Conditional UI Fields**: Hide/show fields (like serial numbers) in the UI based on Item Type `supported_features`.

## Phase 29: Entity Management & Optimization [DONE]
- [x] **Location Dropdown Fix**: Ensure the Asset creation dropdown correctly populates from the Places API.
- [x] **Personnel Edit Fix**: Resolve issues with the Personnel edit page failing to load/save email, phone, and role.
- [x] **View Profile Implementation**: Flesh out the "View Profile" link for Personnel.
- [x] **Inspection Template Polish**: Fix the template edit page and implement a way to assign templates to Item Types.

## Phase 30: System Admin & IAM [DONE]
- [x] User Management (Update Role, Enable/Disable, Delete)
- [x] Global System Settings (Identity, Logistics, Feature Flags)
- [x] Resolved: Catalog "View Details" navigation and inspection sync bugs
