# Rental Management System - Milestone VI Roadmap

This milestone focuses on closing the loop between deployment, to target, to maintenance, to inspection, to re-deployment.

## Phase 23: Entity Management

Goal: Implement first class entities to represent the deployment targets: Companies, Company Contacts, Company Sites (Facilities), Company Locations (Rooms), Company Events (Projects), Event Asset Needs.

### Companies

- [ ] **Company Management**: Implement a UI for creating and managing companies.
- [ ] **Company Contact Management**: Implement a UI for creating and managing company contacts. Contacts represent individuals who may be associated with a company.
- [ ] **Company Site Management**: Implement a UI for creating and managing company site. Sites represent facilities utilized by a company which possess a mailing address.
- [ ] **Company Location Management**: Implement a UI for creating and managing company locations. Locations represent physical spaces within a company site. Locations may be nested to represent subspaces. Locations may have presumed asset needs, which will be applied to an Event Asset Need Calculation automatically (unless overridden), if that location is included in the event.
- [ ] **Company Event Management**: Implement a UI for creating and managing company events. Events have a start and end date, belong to a company, and may reference any number of Contacts, Sites, Locations, and Asset Needs which are disconnected from a given location.
- [ ] **Company Event Asset Need Management**: Implement a UI for creating and managing company asset needs for each event.
- [ ] **Company, Site, Location, Event, and Asset Need API Integration**: Implement a pluggable interface for integration with generic outside REST services.

## Future Plans

- [ ] **Multi-Facility Synchronization**: Global inventory visibility across geographically distributed sites.
- [ ] **AI-Driven Logistics**: Predictive maintenance and automated reordering based on utilization trends.
- [ ] **Third-Party Logistics (3PL) Integration**: Connect with shipping partners for automated dispatch.
- [ ] **MQTT Command Ingest**: Allow MQTT-connected clients to submit `RentAction` requests or control devices directly.
- [ ] **Mobile Technical Persona**: Native-like mobile experience for technicians performing inspections on-site.
- [ ] **Offline Mode**: Support for `RentAction` creation and `Kiosk` scanning in low-connectivity environments.
