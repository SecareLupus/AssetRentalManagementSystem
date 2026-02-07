# FIXMEs for each Page/Interface in our current Built-in UI

- [ ] Dashboard
  - [ ] We currently have 4 scorecards in a line which breaks after 3. We should either break after 2, or increase the number of scorecards to 6, so we have a consistent 2x3 grid.
  - [ ] Stats in the scorecards are currently broken, possibly caused by the GET error listed in Note #1
- [ ] Catalog
  - [ ] The filter button is unused and should be removed.
- [ ] Fleet Reports
  - [ ] Clicking this button in the menu results in a blank page.
- [ ] Settings
  - [ ] The settings menu button currently does nothing
  - [ ] We should flesh out settings, both a user level settings page for this button, as well as managing server level settings in a Server Admin page (not the System Admin page, which manages administrative settings for the Rental Management System itself)
  - [ ] Related to the Server Admin page, we don't currently have any way to manage our IAM or user directory from the UI.
- [ ] Item Types
  - [ ] When creating a new item type, it expects a SKU for fungible items. Does this make sense? What about the kits? We should double check our assumptions about how Item Type skus work.
  - [ ] Supported features checkboxes should affect visibility of those features/fields in the UI, but not affect how they are represented in the database. For example, if we don't support serial numbers for an item type, we still keep an empty field for the serial number in the database, but we should not show the serial number field in the UI.
- [ ] Assets
  - [ ] We should revisit how we manage the identifying codes (skus, serial numbers, etc.) for assets (and item types). I want to make sure we're not conflating component identifiers with company product codes or serial numbers.
  - [ ] Similarly we should consider how to manage serial numbers for internal components and how they get tracked if the internal components get swapped, either during refurbishment or in the field.
  - [ ] Currently cannot create an asset without a location.
  - [ ] Even when there are locations in the backend, the dropdown does not populate with any of them.
- [ ] Locations
  - [ ] We should establish a default location, considered internal (managed by the user's company) for the purpose of tracking all assets not currently deployed.
  - [ ] Any asset not provided a location should default to the unnamed internal location, or the default internal location if one is designated.
  - [ ] This will provide for a simplified workflow for small organizations that don't have a need for managing internal locations.
  - [ ] The Locations page currently reports a GET error (see Note #2)
  - [ ] In the same way that we do not require a defined internal location for simplified workflows, we also should not require a defined company to represent the organization that owns the RMS. Internal locations should not require a company, as they are assumed to be owned by the user's company.
  - [ ] The list of types of places, and their respective subsets of Place data should be configurable at the system settings level.
- [ ] Personnel
  - [ ] After saving a user, and clicking the edit button on their card, the edit page is not populated with email, phone, company, or role. It's unclear whether it's failing to save or failing to load.
  - [ ] View Profile link does nothing.
- [ ] Inspection Templates
  - [ ] Edit page doesn't populate the current state of the template being edited. Data is being saved, as evident by the list of templates, but the edit page is blank.
  - [ ] There is no apparent way to assign inspection templates to item types.

# Notes

#= Note #1 =================================================================#

GET /v1/entities/places

Request
`"No data"`

Response (500)
`"sql: Scan error on column index 8, name \"presumed_demands\": unsupported Scan, storing driver.Value type <nil> into type *json.RawMessage\n"`

#= Note #2 =================================================================#

GET /v1/dashboard/stats
Request

`"No data"`

Response (500)
`"pq: relation \"rent_actions\" does not exist at column 22 (42P01)\n"`

#===========================================================================#