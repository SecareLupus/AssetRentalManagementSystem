# FIXMEs for each Page/Interface in our current Built-in UI

- [ ] Dashboard
  - [ ] Category overview list items not interactive
  - [ ] Category overview view all button not working
- [ ] Catalog
  - [ ] Filters button not working
  - [ ] Add to cart button not working
- [ ] Item Type
  - [ ] Archive button exists on already archived item types
  - [ ] Can't edit "Supported Features" (Let me know if this is necessary given our use case)
  - [ ] Assigned "Supported Features" don't currently appear to do anything
  - [ ] Request Reservation button not working
  - [ ] List of assets of Item Type not interactive
- [ ] Reservation Wizard
  - [ ] Default Dates (Should be today and 7 days from today)
  - [ ] Unclear what ASAP Priority does
  - [ ] Clicking the same item type multiple times should increase quantity rather than adding a duplicate row
  - [ ] Failed to create reservation (See Note #1 Below for API Query Details)
- [ ] Maintenance
  - [ ] Implement an Inspection Template Editor
  - [ ] Failed to submit inspection (See Note #2 Below for API Query Details)
- [ ] Warehouse
  - [ ] Implement a Scannable Tag Editor for mapping aliases to assets
    - [ ] Option for regex matching values in third party scannable tags (eg QR Codes) to extract asset tags or aliases from the scanned value.
  - [ ] Bulk Transactions have no way to accept input (eg start/end dates for checkouts, locations, etc)
- [ ] Intelligence Hub
  - [ ] Critical Shortage value should be editable (probably in the Item Type settings)
  - [ ] Service Forecast Schedule Frequency should be editable (probably in the Item Type settings)
  - [ ] Service Forecast Schedule Snooze option for devices stuck in the field.
  - [ ] Launch Simulator button not working
  - [ ] Heatmap doesn't respect bulk checkouts without schedules
  - [ ] Heamap pagination buttons don't work
- [ ] Simulator
  - [ ] Non-overlapping scenarios do not treat the returned assets as available for later scenarios. (The second scenario will report the asset as unavailable)

# Notes

Note #1: Reservation Creation API Query Details
Endpoint: `POST /v1/rent-actions`

Request:

```
"{
	'requester_ref':'admin',
	'created_by_ref':'admin',
	'priority':'high',
	'start_time':'2026-02-07T13:44:00.000Z',
	'end_time':'2026-02-08T13:44:00.000Z',
	'is_asap':true,
	'description':'Throwing them off a cliff, ...',
	'items' [
		{
			'item_kind':'item_type',
			'item_id':4,
			'requested_quantity':1,
			'name':'Encoder Box'
		}
	],
	'status':'draft'
}"
```

Response:

```
"insert rent_action: pq:
INSERT has more expressions than target columns at position 5:88 (42601)"
```

Note #2: Maintenance API Query Details
Endpoint: `POST /v1/inventory/assets/1/inspections`

Request:

```
"{
	'performed_by':'admin',
	'responses':[]
}"
```

Response:

```
"insert submission: pq:
insert or update on table 'inspection_submissions' violates foreign key constraint 'fk_is_template' (23503)"
```
