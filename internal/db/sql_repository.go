package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/lib/pq"
)

type SqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(db *sql.DB) *SqlRepository {
	return &SqlRepository{db: db}
}

// CreateItemType creates a new item type.
func (r *SqlRepository) CreateItemType(ctx context.Context, it *domain.ItemType) error {
	now := time.Now()
	it.CreatedAt = now
	it.UpdatedAt = now

	query := `INSERT INTO item_types (code, name, kind, is_active, supported_features, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	featuresJSON, _ := json.Marshal(it.SupportedFeatures)
	err := r.db.QueryRowContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, featuresJSON, it.CreatedByUserID, it.UpdatedByUserID, it.SchemaOrg, it.Metadata, it.CreatedAt, it.UpdatedAt,
	).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("create item_type: %w", err)
	}
	return nil
}

// GetItemTypeByID retrieves an item type by its ID.
func (r *SqlRepository) GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error) {
	query := `SELECT id, code, name, kind, is_active, supported_features, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM item_types WHERE id = $1`

	var it domain.ItemType
	var featuresJSON, schemaOrgJSON, metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &featuresJSON, &it.CreatedByUserID, &it.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &it.CreatedAt, &it.UpdatedAt,
	)
	if err == nil {
		json.Unmarshal(featuresJSON, &it.SupportedFeatures)
		it.SchemaOrg = json.RawMessage(schemaOrgJSON)
		it.Metadata = json.RawMessage(metadataJSON)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan item_type: %w", err)
	}
	return &it, nil
}

// ListItemTypes returns item types, optionally including inactive ones.
func (r *SqlRepository) ListItemTypes(ctx context.Context, includeInactive bool) ([]domain.ItemType, error) {
	query := `SELECT id, code, name, kind, is_active, supported_features, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM item_types`
	if !includeInactive {
		query += ` WHERE is_active = TRUE`
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query item_types: %w", err)
	}
	defer rows.Close()

	var results []domain.ItemType
	for rows.Next() {
		var it domain.ItemType
		var featuresJSON, schemaOrgJSON, metadataJSON []byte
		if err := rows.Scan(&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &featuresJSON, &it.CreatedByUserID, &it.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item_type: %w", err)
		}
		json.Unmarshal(featuresJSON, &it.SupportedFeatures)
		it.SchemaOrg = json.RawMessage(schemaOrgJSON)
		it.Metadata = json.RawMessage(metadataJSON)
		results = append(results, it)
	}
	return results, nil
}

// UpdateItemType updates an existing item type.
func (r *SqlRepository) UpdateItemType(ctx context.Context, it *domain.ItemType) error {
	it.UpdatedAt = time.Now()
	query := `UPDATE item_types SET code = $1, name = $2, kind = $3, is_active = $4, supported_features = $5, updated_by_user_id = $6, schema_org = $7, metadata = $8, updated_at = $9
	          WHERE id = $10`

	featuresJSON, _ := json.Marshal(it.SupportedFeatures)
	_, err := r.db.ExecContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, featuresJSON, it.UpdatedByUserID, it.SchemaOrg, it.Metadata, it.UpdatedAt, it.ID,
	)
	if err != nil {
		return fmt.Errorf("update item_type: %w", err)
	}
	return nil
}

// DeleteItemType soft deletes an item type.
func (r *SqlRepository) DeleteItemType(ctx context.Context, id int64) error {
	query := `UPDATE item_types SET is_active = FALSE, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("delete item_type: %w", err)
	}
	return nil
}

// GetAssetByID retrieves a specific asset by its ID.
func (r *SqlRepository) GetAssetByID(ctx context.Context, id int64) (*domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, place_id, location, assigned_to, mesh_node_id, wireguard_hostname, management_url,
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 usage_hours, next_service_hours, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE id = $1`

	var a domain.Asset
	var schemaOrgJSON, metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.PlaceID, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
		&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
		&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get asset by id: %w", err)
	}
	a.SchemaOrg = json.RawMessage(schemaOrgJSON)
	a.Metadata = json.RawMessage(metadataJSON)

	// Phase 28: Extract components from Metadata
	if len(a.Metadata) > 0 {
		var meta struct {
			Components []domain.Component `json:"components"`
		}
		if err := json.Unmarshal(a.Metadata, &meta); err == nil {
			a.Components = meta.Components
		}
	}

	return &a, nil
}

// CreateAsset creates a new asset.
func (r *SqlRepository) CreateAsset(ctx context.Context, a *domain.Asset) error {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	// Phase 28: Default Location Assignment
	if a.PlaceID == nil {
		place, err := r.GetDefaultInternalPlace(ctx)
		if err != nil {
			return fmt.Errorf("get default internal place: %w", err)
		}
		a.PlaceID = &place.ID
	}

	// Phase 28: Status inferred from location
	place, err := r.GetPlace(ctx, *a.PlaceID)
	if err != nil {
		return fmt.Errorf("get place for status inference: %w", err)
	}
	if a.Status == "" {
		if place.IsInternal {
			a.Status = domain.AssetStatusAvailable
		} else {
			a.Status = domain.AssetStatusDeployed
		}
	}

	// Phase 28: Components tracking stored in Metadata
	if len(a.Components) > 0 {
		var meta map[string]interface{}
		if len(a.Metadata) > 0 {
			json.Unmarshal(a.Metadata, &meta)
		} else {
			meta = make(map[string]interface{})
		}
		meta["components"] = a.Components
		a.Metadata, _ = json.Marshal(meta)
	}

	query := `INSERT INTO assets (
		item_type_id, asset_tag, serial_number, status, place_id, location, assigned_to, 
		mesh_node_id, wireguard_hostname, management_url, build_spec_version, provisioning_status, 
		firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
		usage_hours, next_service_hours, schema_org, metadata, created_by_user_id, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24) RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.PlaceID, a.Location, a.AssignedTo,
		a.MeshNodeID, a.WireguardHostname, a.ManagementURL, a.BuildSpecVersion, a.ProvisioningStatus,
		a.FirmwareVersion, a.Hostname, a.RemoteManagementID, a.CurrentBuildSpecID, a.LastInspectionAt,
		a.UsageHours, a.NextServiceHours, a.SchemaOrg, a.Metadata, a.CreatedByUserID, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.ID)
	if err != nil {
		return fmt.Errorf("create asset: %w", err)
	}
	return nil
}

// ListAssets returns all assets.
func (r *SqlRepository) ListAssets(ctx context.Context) ([]domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, place_id, location, assigned_to, mesh_node_id, wireguard_hostname, management_url, 
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 usage_hours, next_service_hours, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM assets`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	var results []domain.Asset
	for rows.Next() {
		var a domain.Asset
		var schemaOrgJSON, metadataJSON []byte
		if err := rows.Scan(
			&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.PlaceID, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
			&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
			&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		a.SchemaOrg = json.RawMessage(schemaOrgJSON)
		a.Metadata = json.RawMessage(metadataJSON)

		// Phase 28: Extract components
		if len(a.Metadata) > 0 {
			var meta struct {
				Components []domain.Component `json:"components"`
			}
			if err := json.Unmarshal(a.Metadata, &meta); err == nil {
				a.Components = meta.Components
			}
		}

		results = append(results, a)
	}
	return results, nil
}

// ListAssetsByItemType returns assets belonging to a specific item type.
func (r *SqlRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, place_id, location, assigned_to, mesh_node_id, wireguard_hostname, management_url, 
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 usage_hours, next_service_hours, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE item_type_id = $1`

	rows, err := r.db.QueryContext(ctx, query, itemTypeID)
	if err != nil {
		return nil, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	var results []domain.Asset
	for rows.Next() {
		var a domain.Asset
		var schemaOrgJSON, metadataJSON []byte
		if err := rows.Scan(
			&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.PlaceID, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
			&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
			&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		a.SchemaOrg = json.RawMessage(schemaOrgJSON)
		a.Metadata = json.RawMessage(metadataJSON)

		// Phase 28: Extract components
		if len(a.Metadata) > 0 {
			var meta struct {
				Components []domain.Component `json:"components"`
			}
			if err := json.Unmarshal(a.Metadata, &meta); err == nil {
				a.Components = meta.Components
			}
		}

		results = append(results, a)
	}
	return results, nil
}

// UpdateAsset updates an existing asset.
func (r *SqlRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	a.UpdatedAt = time.Now()
	query := `UPDATE assets SET 
		item_type_id = $1, asset_tag = $2, serial_number = $3, status = $4, 
		place_id = $5, location = $6, assigned_to = $7, mesh_node_id = $8, wireguard_hostname = $9,
		management_url = $10, build_spec_version = $11, provisioning_status = $12, firmware_version = $13,
		hostname = $14, remote_management_id = $15, current_build_spec_id = $16, last_inspection_at = $17,
		usage_hours = $18, next_service_hours = $19, updated_by_user_id = $20, schema_org = $21, 
		metadata = $22, updated_at = $23
		WHERE id = $24`

	// Phase 28: Components tracking stored in Metadata
	if len(a.Components) > 0 {
		var meta map[string]interface{}
		if len(a.Metadata) > 0 {
			json.Unmarshal(a.Metadata, &meta)
		} else {
			meta = make(map[string]interface{})
		}
		meta["components"] = a.Components
		a.Metadata, _ = json.Marshal(meta)
	}

	_, err := r.db.ExecContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.PlaceID, a.Location, a.AssignedTo,
		a.MeshNodeID, a.WireguardHostname, a.ManagementURL, a.BuildSpecVersion, a.ProvisioningStatus,
		a.FirmwareVersion, a.Hostname, a.RemoteManagementID, a.CurrentBuildSpecID, a.LastInspectionAt,
		a.UsageHours, a.NextServiceHours, a.UpdatedByUserID, a.SchemaOrg, a.Metadata, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return fmt.Errorf("update asset: %w", err)
	}
	return nil
}

// UpdateAssetStatus updates the status of an asset along with optional metadata and location.
func (r *SqlRepository) UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus, placeID *int64, location *string, metadata json.RawMessage) error {
	query := `UPDATE assets SET status = $1, updated_at = $2`
	args := []interface{}{status, time.Now()}
	argCount := 3

	if placeID != nil {
		query += fmt.Sprintf(", place_id = $%d", argCount)
		args = append(args, *placeID)
		argCount++
	}

	if location != nil {
		query += fmt.Sprintf(", location = $%d", argCount)
		args = append(args, *location)
		argCount++
	}

	if metadata != nil {
		query += fmt.Sprintf(", metadata = $%d", argCount)
		args = append(args, metadata)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update asset status: %w", err)
	}
	return nil
}

// RecallAssetsByItemType moves all deployed/available assets of a type into 'recalled' status.
func (r *SqlRepository) RecallAssetsByItemType(ctx context.Context, itemTypeID int64) error {
	query := `UPDATE assets SET status = 'recalled', updated_at = $1 
	          WHERE item_type_id = $2 AND (status = 'available' OR status = 'deployed')`
	_, err := r.db.ExecContext(ctx, query, time.Now(), itemTypeID)
	if err != nil {
		return fmt.Errorf("bulk recall assets: %w", err)
	}
	return nil
}

// DeleteAsset deletes an asset (permanent).
func (r *SqlRepository) DeleteAsset(ctx context.Context, id int64) error {
	query := `DELETE FROM assets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete asset: %w", err)
	}
	return nil
}

// Logistics

// CreateRentalReservation creates a new rental reservation and its associated demands.
func (r *SqlRepository) CreateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	rr.CreatedAt = now
	rr.UpdatedAt = now
	if rr.BookingTime.IsZero() {
		rr.BookingTime = now
	}

	query := `INSERT INTO rental_reservations (
		reservation_name, reservation_status, under_name_id, booking_time, 
		start_time, end_time, provider_id, metadata, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		rr.ReservationName, rr.ReservationStatus, rr.UnderNameID, rr.BookingTime,
		rr.StartTime, rr.EndTime, rr.ProviderID, rr.Metadata, rr.CreatedAt, rr.UpdatedAt,
	).Scan(&rr.ID)
	if err != nil {
		return fmt.Errorf("insert rental_reservation: %w", err)
	}

	for i := range rr.Demands {
		d := &rr.Demands[i]
		d.ReservationID = rr.ID
		d.CreatedAt = now
		d.UpdatedAt = now
		dQuery := `INSERT INTO demands (
			reservation_id, event_id, item_kind, item_id, requested_quantity, 
			business_function, eligible_duration, place_id, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

		err = tx.QueryRowContext(ctx, dQuery,
			d.ReservationID, d.EventID, d.ItemKind, d.ItemID, d.Quantity,
			d.BusinessFunction, d.EligibleDuration, d.PlaceID, d.Metadata, d.CreatedAt, d.UpdatedAt,
		).Scan(&d.ID)
		if err != nil {
			return fmt.Errorf("insert demand: %w", err)
		}
	}

	return tx.Commit()
}

// GetRentalReservationByID retrieves a reservation by its ID, including its demands.
func (r *SqlRepository) GetRentalReservationByID(ctx context.Context, id int64) (*domain.RentalReservation, error) {
	query := `SELECT id, reservation_name, reservation_status, under_name_id, booking_time, 
	                 start_time, end_time, provider_id, metadata, created_at, updated_at 
	          FROM rental_reservations WHERE id = $1`

	var rr domain.RentalReservation
	var metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rr.ID, &rr.ReservationName, &rr.ReservationStatus, &rr.UnderNameID, &rr.BookingTime,
		&rr.StartTime, &rr.EndTime, &rr.ProviderID, &metadataJSON, &rr.CreatedAt, &rr.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan rental_reservation: %w", err)
	}
	rr.Metadata = json.RawMessage(metadataJSON)

	demandQuery := `SELECT id, reservation_id, event_id, item_kind, item_id, requested_quantity, 
	                       business_function, eligible_duration, place_id, metadata, created_at, updated_at 
	                FROM demands WHERE reservation_id = $1`

	rows, err := r.db.QueryContext(ctx, demandQuery, id)
	if err != nil {
		return nil, fmt.Errorf("query demands: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d domain.Demand
		var dMetadataJSON []byte
		if err := rows.Scan(
			&d.ID, &d.ReservationID, &d.EventID, &d.ItemKind, &d.ItemID, &d.Quantity,
			&d.BusinessFunction, &d.EligibleDuration, &d.PlaceID, &dMetadataJSON, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan demand: %w", err)
		}
		d.Metadata = json.RawMessage(dMetadataJSON)
		rr.Demands = append(rr.Demands, d)
	}

	return &rr, nil
}

// ListRentalReservations returns all rental reservations.
func (r *SqlRepository) ListRentalReservations(ctx context.Context) ([]domain.RentalReservation, error) {
	query := `SELECT id, reservation_name, reservation_status, under_name_id, booking_time, 
	                 start_time, end_time, provider_id, metadata, created_at, updated_at 
	          FROM rental_reservations ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rental_reservations: %w", err)
	}
	defer rows.Close()

	var results []domain.RentalReservation
	for rows.Next() {
		var rr domain.RentalReservation
		var metadataJSON []byte
		err := rows.Scan(
			&rr.ID, &rr.ReservationName, &rr.ReservationStatus, &rr.UnderNameID, &rr.BookingTime,
			&rr.StartTime, &rr.EndTime, &rr.ProviderID, &metadataJSON, &rr.CreatedAt, &rr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rental_reservation: %w", err)
		}
		rr.Metadata = json.RawMessage(metadataJSON)
		results = append(results, rr)
	}
	return results, nil
}

// UpdateRentalReservation updates an existing reservation.
func (r *SqlRepository) UpdateRentalReservation(ctx context.Context, rr *domain.RentalReservation) error {
	rr.UpdatedAt = time.Now()
	query := `UPDATE rental_reservations SET 
		reservation_name = $1, reservation_status = $2, under_name_id = $3, booking_time = $4, 
		start_time = $5, end_time = $6, provider_id = $7, metadata = $8, updated_at = $9 
		WHERE id = $10`

	_, err := r.db.ExecContext(ctx, query,
		rr.ReservationName, rr.ReservationStatus, rr.UnderNameID, rr.BookingTime,
		rr.StartTime, rr.EndTime, rr.ProviderID, rr.Metadata, rr.UpdatedAt, rr.ID,
	)
	if err != nil {
		return fmt.Errorf("update rental_reservation: %w", err)
	}
	return nil
}

// UpdateRentalReservationStatus updates only the status for a reservation.
func (r *SqlRepository) UpdateRentalReservationStatus(ctx context.Context, id int64, status domain.RentalReservationStatus) error {
	query := "UPDATE rental_reservations SET reservation_status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// Demand Methods

func (r *SqlRepository) CreateDemand(ctx context.Context, d *domain.Demand) error {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	query := `INSERT INTO demands (
		reservation_id, event_id, item_kind, item_id, requested_quantity, 
		business_function, eligible_duration, place_id, metadata, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		d.ReservationID, d.EventID, d.ItemKind, d.ItemID, d.Quantity,
		d.BusinessFunction, d.EligibleDuration, d.PlaceID, d.Metadata, d.CreatedAt, d.UpdatedAt,
	).Scan(&d.ID)
	return err
}

func (r *SqlRepository) ListDemandsByReservation(ctx context.Context, reservationID int64) ([]domain.Demand, error) {
	query := `SELECT id, reservation_id, event_id, item_kind, item_id, requested_quantity, 
	                 business_function, eligible_duration, place_id, metadata, created_at, updated_at 
	          FROM demands WHERE reservation_id = $1`
	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Demand
	for rows.Next() {
		var d domain.Demand
		var metadataJSON []byte
		if err := rows.Scan(&d.ID, &d.ReservationID, &d.EventID, &d.ItemKind, &d.ItemID, &d.Quantity, &d.BusinessFunction, &d.EligibleDuration, &d.PlaceID, &metadataJSON, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		d.Metadata = json.RawMessage(metadataJSON)
		results = append(results, d)
	}
	return results, nil
}

func (r *SqlRepository) ListDemandsByEvent(ctx context.Context, eventID int64) ([]domain.Demand, error) {
	query := `SELECT id, reservation_id, event_id, item_kind, item_id, requested_quantity, 
	                 business_function, eligible_duration, place_id, metadata, created_at, updated_at 
	          FROM demands WHERE event_id = $1`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Demand
	for rows.Next() {
		var d domain.Demand
		var metadataJSON []byte
		if err := rows.Scan(&d.ID, &d.ReservationID, &d.EventID, &d.ItemKind, &d.ItemID, &d.Quantity, &d.BusinessFunction, &d.EligibleDuration, &d.PlaceID, &metadataJSON, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		d.Metadata = json.RawMessage(metadataJSON)
		results = append(results, d)
	}
	return results, nil
}

func (r *SqlRepository) UpdateDemand(ctx context.Context, d *domain.Demand) error {
	d.UpdatedAt = time.Now()
	query := `UPDATE demands SET 
		reservation_id = $1, event_id = $2, item_kind = $3, item_id = $4, requested_quantity = $5, 
		business_function = $6, eligible_duration = $7, place_id = $8, metadata = $9, updated_at = $10 
		WHERE id = $11`
	_, err := r.db.ExecContext(ctx, query, d.ReservationID, d.EventID, d.ItemKind, d.ItemID, d.Quantity, d.BusinessFunction, d.EligibleDuration, d.PlaceID, d.Metadata, d.UpdatedAt, d.ID)
	return err
}

func (r *SqlRepository) DeleteDemand(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM demands WHERE id = $1", id)
	return err
}

// Action Methods

func (r *SqlRepository) CreateCheckOutAction(ctx context.Context, co *domain.CheckOutAction) error {
	query := `INSERT INTO check_out_actions (
		reservation_id, asset_id, agent_id, recipient_id, start_time, 
		from_location_id, to_location_id, action_status, metadata
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return r.db.QueryRowContext(ctx, query, co.ReservationID, co.AssetID, co.AgentID, co.RecipientID, co.StartTime, co.FromLocation, co.ToLocation, co.Status, co.Metadata).Scan(&co.ID)
}

func (r *SqlRepository) CreateReturnAction(ctx context.Context, ra *domain.ReturnAction) error {
	query := `INSERT INTO return_actions (
		reservation_id, asset_id, agent_id, start_time, 
		from_location_id, to_location_id, action_status, metadata
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRowContext(ctx, query, ra.ReservationID, ra.AssetID, ra.AgentID, ra.StartTime, ra.FromLocation, ra.ToLocation, ra.Status, ra.Metadata).Scan(&ra.ID)
}

func (r *SqlRepository) ListCheckOutActions(ctx context.Context, reservationID int64) ([]domain.CheckOutAction, error) {
	query := `SELECT id, reservation_id, asset_id, agent_id, recipient_id, start_time, 
	                 from_location_id, to_location_id, action_status, metadata 
	          FROM check_out_actions WHERE reservation_id = $1`
	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.CheckOutAction
	for rows.Next() {
		var co domain.CheckOutAction
		var metadataJSON []byte
		if err := rows.Scan(&co.ID, &co.ReservationID, &co.AssetID, &co.AgentID, &co.RecipientID, &co.StartTime, &co.FromLocation, &co.ToLocation, &co.Status, &metadataJSON); err != nil {
			return nil, err
		}
		co.Metadata = json.RawMessage(metadataJSON)
		results = append(results, co)
	}
	return results, nil
}

func (r *SqlRepository) ListReturnActions(ctx context.Context, reservationID int64) ([]domain.ReturnAction, error) {
	query := `SELECT id, reservation_id, asset_id, agent_id, start_time, 
	                 from_location_id, to_location_id, action_status, metadata 
	          FROM return_actions WHERE reservation_id = $1`
	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.ReturnAction
	for rows.Next() {
		var ra domain.ReturnAction
		var metadataJSON []byte
		if err := rows.Scan(&ra.ID, &ra.ReservationID, &ra.AssetID, &ra.AgentID, &ra.StartTime, &ra.FromLocation, &ra.ToLocation, &ra.Status, &metadataJSON); err != nil {
			return nil, err
		}
		ra.Metadata = json.RawMessage(metadataJSON)
		results = append(results, ra)
	}
	return results, nil
}

// GetRentalFulfillmentStatus calculates the delta between demands and actual movements.
func (r *SqlRepository) GetRentalFulfillmentStatus(ctx context.Context, reservationID int64) (*domain.RentalFulfillmentStatus, error) {
	demands, err := r.ListDemandsByReservation(ctx, reservationID)
	if err != nil {
		return nil, err
	}

	// Map to track fulfillment per demand
	// For simplicity, we assume one demand per (item_kind, item_id)
	lineMap := make(map[string]*domain.FulfillmentLine)
	for _, d := range demands {
		key := fmt.Sprintf("%s:%d", d.ItemKind, d.ItemID)
		lineMap[key] = &domain.FulfillmentLine{
			DemandID:          d.ID,
			ItemKind:          d.ItemKind,
			ItemID:            d.ItemID,
			RequestedQuantity: d.Quantity,
		}
	}

	// Better approach: use a combined query or join
	coQuery := `
		SELECT a.item_type_id, COUNT(*) 
		FROM check_out_actions co
		JOIN assets a ON co.asset_id = a.id
		WHERE co.reservation_id = $1 AND co.action_status = 'Completed'
		GROUP BY a.item_type_id
	`
	rows, err := r.db.QueryContext(ctx, coQuery, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var itemID int64
		var count int
		if err := rows.Scan(&itemID, &count); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("item_type:%d", itemID)
		if line, ok := lineMap[key]; ok {
			line.FulfilledQuantity = count
		}
	}

	retQuery := `
		SELECT a.item_type_id, COUNT(*) 
		FROM return_actions ret
		JOIN assets a ON ret.asset_id = a.id
		WHERE ret.reservation_id = $1 AND ret.action_status = 'Completed'
		GROUP BY a.item_type_id
	`
	rows2, err := r.db.QueryContext(ctx, retQuery, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var itemID int64
		var count int
		if err := rows2.Scan(&itemID, &count); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("item_type:%d", itemID)
		if line, ok := lineMap[key]; ok {
			line.ReturnedQuantity = count
		}
	}

	var status domain.RentalFulfillmentStatus
	status.ReservationID = reservationID

	anyFulfilled := false
	allFulfilled := true
	for _, line := range lineMap {
		line.RemainingQuantity = line.RequestedQuantity - line.FulfilledQuantity
		if line.FulfilledQuantity > 0 {
			anyFulfilled = true
		}
		if line.FulfilledQuantity < line.RequestedQuantity {
			allFulfilled = false
		}
		status.Lines = append(status.Lines, *line)
	}

	if allFulfilled && len(lineMap) > 0 {
		status.Status = string(domain.ReservationStatusFulfilled)
	} else if anyFulfilled {
		status.Status = string(domain.ReservationStatusPartiallyFulfilled)
	} else {
		status.Status = string(domain.ReservationStatusConfirmed)
	}

	return &status, nil
}

// BatchCheckOut dispatches multiple assets at once and updates reservation status.
func (r *SqlRepository) BatchCheckOut(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64, fromLocationID, toLocationID *int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	for _, assetID := range assetIDs {
		// 1. Create CheckOutAction
		coQuery := `INSERT INTO check_out_actions (reservation_id, asset_id, agent_id, start_time, from_location_id, to_location_id, action_status)
		            VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = tx.ExecContext(ctx, coQuery, reservationID, assetID, agentID, now, fromLocationID, toLocationID, "Completed")
		if err != nil {
			return fmt.Errorf("checkout %d: %w", assetID, err)
		}

		// 2. Update Asset
		assetQuery := `UPDATE assets SET status = 'deployed', place_id = $1, updated_at = $2 WHERE id = $3`
		_, err = tx.ExecContext(ctx, assetQuery, toLocationID, now, assetID)
		if err != nil {
			return fmt.Errorf("update asset %d: %w", assetID, err)
		}
	}

	// 3. Update Reservation Status based on new fulfillment state
	// (Normally we'd re-calculate here, but for simplicity we'll just set to PartiallyFulfilled if any movement happens)
	// Improved: we can fetch the status and update it.

	err = tx.Commit()
	if err != nil {
		return err
	}

	// Re-evaluate overall status after commit
	fStatus, err := r.GetRentalFulfillmentStatus(ctx, reservationID)
	if err == nil {
		r.UpdateRentalReservationStatus(ctx, reservationID, domain.RentalReservationStatus(fStatus.Status))
	}

	return nil
}

// BatchReturn returns multiple assets at once.
func (r *SqlRepository) BatchReturn(ctx context.Context, reservationID int64, assetIDs []int64, agentID int64, toLocationID *int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	for _, assetID := range assetIDs {
		// 1. Create ReturnAction
		raQuery := `INSERT INTO return_actions (reservation_id, asset_id, agent_id, start_time, to_location_id, action_status)
		            VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, raQuery, reservationID, assetID, agentID, now, toLocationID, "Completed")
		if err != nil {
			return fmt.Errorf("return %d: %w", assetID, err)
		}

		// 2. Update Asset to available
		assetQuery := `UPDATE assets SET status = 'available', place_id = $1, updated_at = $2 WHERE id = $3`
		_, err = tx.ExecContext(ctx, assetQuery, toLocationID, now, assetID)
		if err != nil {
			return fmt.Errorf("update asset %d: %w", assetID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetAvailableQuantity calculates the available inventory for an item type in a given time window.
func (r *SqlRepository) GetAvailableQuantity(ctx context.Context, itemTypeID int64, startTime, endTime time.Time) (int, error) {
	// 1. Get total assets for this item type (excluding retired)
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM assets WHERE item_type_id = $1 AND status != 'retired'", itemTypeID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count assets: %w", err)
	}

	// 2. Subtract quantities from overlapping CONFIRMED reservations
	// Schema.org: ReservationConfirmed
	query := `
		SELECT COALESCE(SUM(d.requested_quantity), 0)
		FROM demands d
		JOIN rental_reservations rr ON d.reservation_id = rr.id
		WHERE d.item_kind = 'item_type' 
		  AND d.item_id = $1
		  AND rr.reservation_status = 'ReservationConfirmed'
		  AND rr.start_time < $3
		  AND rr.end_time > $2
	`
	var reserved int
	err = r.db.QueryRowContext(ctx, query, itemTypeID, startTime, endTime).Scan(&reserved)
	if err != nil {
		return 0, fmt.Errorf("sum reserved quantity: %w", err)
	}

	// 3. Subtract Ad-Hoc Usage
	queryAdHoc := `
		SELECT COUNT(*) FROM assets 
		WHERE item_type_id = $1 
		  AND status IN ('deployed', 'maintenance')
		  AND (metadata->>'estimated_return_at' IS NULL OR (metadata->>'estimated_return_at')::timestamp > $2)
	`
	var adHoc int
	err = r.db.QueryRowContext(ctx, queryAdHoc, itemTypeID, startTime).Scan(&adHoc)
	if err != nil {
		return 0, fmt.Errorf("count ad-hoc usage: %w", err)
	}

	return total - reserved - adHoc, nil
}

// AddMaintenanceLog records a new maintenance activity.
func (r *SqlRepository) AddMaintenanceLog(ctx context.Context, ml *domain.MaintenanceLog) error {
	ml.CreatedAt = time.Now()
	query := `INSERT INTO maintenance_logs (asset_id, action_type, notes, performed_by, test_bits, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		ml.AssetID, ml.ActionType, ml.Notes, ml.PerformedBy, ml.TestBits, ml.CreatedAt,
	).Scan(&ml.ID)
	if err != nil {
		return fmt.Errorf("add maintenance log: %w", err)
	}
	return nil
}

// ListMaintenanceLogs retrieves history for a specific asset.
func (r *SqlRepository) ListMaintenanceLogs(ctx context.Context, assetID int64) ([]domain.MaintenanceLog, error) {
	query := `SELECT id, asset_id, action_type, notes, performed_by, test_bits, created_at
	          FROM maintenance_logs WHERE asset_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("query maintenance_logs: %w", err)
	}
	defer rows.Close()

	var results []domain.MaintenanceLog
	for rows.Next() {
		var ml domain.MaintenanceLog
		if err := rows.Scan(&ml.ID, &ml.AssetID, &ml.ActionType, &ml.Notes, &ml.PerformedBy, &ml.TestBits, &ml.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maintenance log: %w", err)
		}
		results = append(results, ml)
	}
	return results, nil
}

// CreateInspectionTemplate creates a new inspection template with associated fields.
func (r *SqlRepository) CreateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	it.CreatedAt = now
	it.UpdatedAt = now

	err = tx.QueryRowContext(ctx, "INSERT INTO inspection_templates (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		it.Name, it.Description, it.CreatedAt, it.UpdatedAt,
	).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("insert template: %w", err)
	}

	for i := range it.Fields {
		f := &it.Fields[i]
		f.TemplateID = it.ID
		err = tx.QueryRowContext(ctx, "INSERT INTO inspection_fields (template_id, label, field_type, required, display_order) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			f.TemplateID, f.Label, f.Type, f.Required, f.DisplayOrder,
		).Scan(&f.ID)
		if err != nil {
			return fmt.Errorf("insert field: %w", err)
		}
	}

	return tx.Commit()
}

// ListInspectionTemplates returns all inspection templates.
func (r *SqlRepository) ListInspectionTemplates(ctx context.Context) ([]domain.InspectionTemplate, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM inspection_templates ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all templates: %w", err)
	}
	defer rows.Close()

	var results []domain.InspectionTemplate
	for rows.Next() {
		var it domain.InspectionTemplate
		if err := rows.Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		results = append(results, it)
	}
	return results, nil
}

// UpdateInspectionTemplate updates a template and its fields.
func (r *SqlRepository) UpdateInspectionTemplate(ctx context.Context, it *domain.InspectionTemplate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	it.UpdatedAt = time.Now()
	_, err = tx.ExecContext(ctx, "UPDATE inspection_templates SET name = $1, description = $2, updated_at = $3 WHERE id = $4",
		it.Name, it.Description, it.UpdatedAt, it.ID,
	)
	if err != nil {
		return fmt.Errorf("update template info: %w", err)
	}

	// Simplest approach for fields: delete and re-insert
	_, err = tx.ExecContext(ctx, "DELETE FROM inspection_fields WHERE template_id = $1", it.ID)
	if err != nil {
		return fmt.Errorf("clear old fields: %w", err)
	}

	for i := range it.Fields {
		f := &it.Fields[i]
		f.TemplateID = it.ID
		err = tx.QueryRowContext(ctx, "INSERT INTO inspection_fields (template_id, label, field_type, required, display_order) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			f.TemplateID, f.Label, f.Type, f.Required, f.DisplayOrder,
		).Scan(&f.ID)
		if err != nil {
			return fmt.Errorf("re-insert field: %w", err)
		}
	}

	return tx.Commit()
}

// DeleteInspectionTemplate removes a template.
func (r *SqlRepository) DeleteInspectionTemplate(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM inspection_templates WHERE id = $1", id)
	return err
}

// SetItemTypeInspections syncs the assignment of templates to an item type.
func (r *SqlRepository) SetItemTypeInspections(ctx context.Context, itemTypeID int64, templateIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM item_type_inspections WHERE item_type_id = $1", itemTypeID)
	if err != nil {
		return fmt.Errorf("clear old assignments: %w", err)
	}

	for _, tid := range templateIDs {
		_, err = tx.ExecContext(ctx, "INSERT INTO item_type_inspections (item_type_id, template_id) VALUES ($1, $2)", itemTypeID, tid)
		if err != nil {
			return fmt.Errorf("insert assignment (it:%d, t:%d): %w", itemTypeID, tid, err)
		}
	}

	return tx.Commit()
}

// GetInspectionTemplate retrieves a single template with its fields.
func (r *SqlRepository) GetInspectionTemplate(ctx context.Context, id int64) (*domain.InspectionTemplate, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM inspection_templates WHERE id = $1`
	var it domain.InspectionTemplate
	err := r.db.QueryRowContext(ctx, query, id).Scan(&it.ID, &it.Name, &it.Description, &it.CreatedAt, &it.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query template: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id, template_id, label, field_type, required, display_order FROM inspection_fields WHERE template_id = $1 ORDER BY display_order", it.ID)
	if err != nil {
		return nil, fmt.Errorf("query fields: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f domain.InspectionField
		if err := rows.Scan(&f.ID, &f.TemplateID, &f.Label, &f.Type, &f.Required, &f.DisplayOrder); err != nil {
			return nil, err
		}
		it.Fields = append(it.Fields, f)
	}
	return &it, nil
}

// GetInspectionTemplatesForItemType retrieves templates assigned to a specific category.
func (r *SqlRepository) GetInspectionTemplatesForItemType(ctx context.Context, itemTypeID int64) ([]domain.InspectionTemplate, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_at, t.updated_at
		FROM inspection_templates t
		JOIN item_type_inspections iti ON t.id = iti.template_id
		WHERE iti.item_type_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, itemTypeID)
	if err != nil {
		return nil, fmt.Errorf("query templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.InspectionTemplate
	for rows.Next() {
		var t domain.InspectionTemplate
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}

		// Fetch fields for each template
		fRows, err := r.db.QueryContext(ctx, "SELECT id, template_id, label, field_type, required, display_order FROM inspection_fields WHERE template_id = $1 ORDER BY display_order", t.ID)
		if err != nil {
			return nil, err
		}
		defer fRows.Close()

		for fRows.Next() {
			var f domain.InspectionField
			if err := fRows.Scan(&f.ID, &f.TemplateID, &f.Label, &f.Type, &f.Required, &f.DisplayOrder); err != nil {
				return nil, err
			}
			t.Fields = append(t.Fields, f)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

// CreateInspectionSubmission records a new inspection result.
func (r *SqlRepository) CreateInspectionSubmission(ctx context.Context, is *domain.InspectionSubmission) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	is.CreatedAt = time.Now()
	err = tx.QueryRowContext(ctx, "INSERT INTO inspection_submissions (asset_id, template_id, performed_by, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		is.AssetID, is.TemplateID, is.PerformedBy, is.CreatedAt,
	).Scan(&is.ID)
	if err != nil {
		return fmt.Errorf("insert submission: %w", err)
	}

	for i := range is.Responses {
		resp := &is.Responses[i]
		resp.SubmissionID = is.ID
		_, err = tx.ExecContext(ctx, "INSERT INTO inspection_responses (submission_id, field_id, response_value) VALUES ($1, $2, $3)",
			resp.SubmissionID, resp.FieldID, resp.Value,
		)
		if err != nil {
			return fmt.Errorf("insert response: %w", err)
		}
	}

	// Update the asset's last inspection timestamp
	_, err = tx.ExecContext(ctx, "UPDATE assets SET last_inspection_at = $1 WHERE id = $2", is.CreatedAt, is.AssetID)
	if err != nil {
		return fmt.Errorf("update asset last_inspection_at: %w", err)
	}

	return tx.Commit()
}

// Build Spec Management

func (r *SqlRepository) CreateBuildSpec(ctx context.Context, bs *domain.BuildSpec) error {
	now := time.Now()
	bs.CreatedAt = now
	bs.UpdatedAt = now

	query := `INSERT INTO build_specs (version, hardware_config, software_config, firmware_url, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		bs.Version, bs.HardwareConfig, bs.SoftwareConfig, bs.FirmwareURL, bs.Metadata, bs.CreatedAt, bs.UpdatedAt,
	).Scan(&bs.ID)
	if err != nil {
		return fmt.Errorf("create build_spec: %w", err)
	}
	return nil
}

func (r *SqlRepository) GetBuildSpecByID(ctx context.Context, id int64) (*domain.BuildSpec, error) {
	query := `SELECT id, version, hardware_config, software_config, firmware_url, metadata, created_at, updated_at 
	          FROM build_specs WHERE id = $1`

	var bs domain.BuildSpec
	var metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bs.ID, &bs.Version, &bs.HardwareConfig, &bs.SoftwareConfig, &bs.FirmwareURL, &metadataJSON, &bs.CreatedAt, &bs.UpdatedAt,
	)
	if err == nil {
		bs.Metadata = json.RawMessage(metadataJSON)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get build_spec: %w", err)
	}
	return &bs, nil
}

func (r *SqlRepository) ListBuildSpecs(ctx context.Context) ([]domain.BuildSpec, error) {
	query := `SELECT id, version, hardware_config, software_config, firmware_url, metadata, created_at, updated_at FROM build_specs`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.BuildSpec
	for rows.Next() {
		var bs domain.BuildSpec
		var metadataJSON []byte
		if err := rows.Scan(&bs.ID, &bs.Version, &bs.HardwareConfig, &bs.SoftwareConfig, &bs.FirmwareURL, &metadataJSON, &bs.CreatedAt, &bs.UpdatedAt); err != nil {
			return nil, err
		}
		bs.Metadata = json.RawMessage(metadataJSON)
		results = append(results, bs)
	}
	return results, nil
}

// Provisioning Workflow

func (r *SqlRepository) StartProvisioning(ctx context.Context, assetID int64, buildSpecID int64, performedBy string) (*domain.ProvisionAction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update Asset status
	_, err = tx.ExecContext(ctx, "UPDATE assets SET status = 'maintenance', provisioning_status = 'flashing', current_build_spec_id = $1 WHERE id = $2", buildSpecID, assetID)
	if err != nil {
		return nil, fmt.Errorf("update asset for provisioning: %w", err)
	}

	// Create ProvisionAction log
	pa := &domain.ProvisionAction{
		AssetID:     assetID,
		BuildSpecID: &buildSpecID,
		Status:      domain.ProvisionStarted,
		PerformedBy: performedBy,
		CreatedAt:   time.Now(),
	}

	query := `INSERT INTO provision_actions (asset_id, build_spec_id, status, performed_by, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.QueryRowContext(ctx, query, pa.AssetID, pa.BuildSpecID, pa.Status, pa.PerformedBy, pa.CreatedAt).Scan(&pa.ID)
	if err != nil {
		return nil, fmt.Errorf("create provision_action: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pa, nil
}

func (r *SqlRepository) CompleteProvisioning(ctx context.Context, actionID int64, notes string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	var assetID int64
	err = tx.QueryRowContext(ctx, "UPDATE provision_actions SET status = 'completed', notes = $1, completed_at = $2 WHERE id = $3 RETURNING asset_id",
		notes, now, actionID,
	).Scan(&assetID)
	if err != nil {
		return fmt.Errorf("complete provision_action: %w", err)
	}

	// Set asset to Ready
	_, err = tx.ExecContext(ctx, "UPDATE assets SET status = 'available', provisioning_status = 'ready' WHERE id = $1", assetID)
	if err != nil {
		return fmt.Errorf("set asset to ready: %w", err)
	}

	return tx.Commit()
}

// User Management

func (r *SqlRepository) CreateUser(ctx context.Context, u *domain.User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := `INSERT INTO users (username, email, password_hash, role, is_enabled, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		u.Username, u.Email, u.PasswordHash, u.Role, u.IsEnabled, u.CreatedAt, u.UpdatedAt,
	).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *SqlRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, is_enabled, last_login_at, created_at, updated_at 
	          FROM users WHERE id = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.IsEnabled, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}

func (r *SqlRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, role, is_enabled, last_login_at, created_at, updated_at 
	          FROM users WHERE username = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.IsEnabled, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &u, nil
}

func (r *SqlRepository) UpdateUser(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now()
	query := `UPDATE users SET email = $1, role = $2, is_enabled = $3, last_login_at = $4, updated_at = $5 WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query, u.Email, u.Role, u.IsEnabled, u.LastLoginAt, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Outbox Implementation

func (r *SqlRepository) AppendEvent(ctx context.Context, tx *sql.Tx, event *domain.OutboxEvent) error {
	query := `INSERT INTO outbox_events (event_type, payload, status, created_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`

	now := time.Now()
	event.CreatedAt = now
	event.Status = domain.OutboxPending

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, event.Type, event.Payload, event.Status, event.CreatedAt).Scan(&event.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, event.Type, event.Payload, event.Status, event.CreatedAt).Scan(&event.ID)
	}

	if err != nil {
		return fmt.Errorf("append outbox event: %w", err)
	}
	return nil
}

func (r *SqlRepository) GetPendingEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	query := `SELECT id, event_type, payload, status, error_message, retry_count, created_at, processed_at
	          FROM outbox_events WHERE status = 'pending' ORDER BY created_at ASC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var e domain.OutboxEvent
		if err := rows.Scan(&e.ID, &e.Type, &e.Payload, &e.Status, &e.ErrorMessage, &e.RetryCount, &e.CreatedAt, &e.ProcessedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *SqlRepository) MarkEventProcessed(ctx context.Context, id int64) error {
	query := `UPDATE outbox_events SET status = 'processed', processed_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *SqlRepository) MarkEventFailed(ctx context.Context, id int64, errMessage string) error {
	query := `UPDATE outbox_events SET status = 'failed', error_message = $1, retry_count = retry_count + 1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, errMessage, id)
	return err
}

func (r *SqlRepository) ListWebhooks(ctx context.Context) ([]domain.WebhookConfig, error) {
	query := `SELECT id, url, secret, enabled_events, is_active, created_at, updated_at FROM webhooks WHERE is_active = TRUE`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.WebhookConfig
	for rows.Next() {
		var w domain.WebhookConfig
		var eventsJSON []byte
		if err := rows.Scan(&w.ID, &w.URL, &w.Secret, &eventsJSON, &w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(eventsJSON, &w.Events)
		results = append(results, w)
	}
	return results, nil
}

// GetAvailabilityTimeline returns availability data points over a range of dates.
func (r *SqlRepository) GetAvailabilityTimeline(ctx context.Context, itemTypeID int64, start, end time.Time) ([]domain.AvailabilityPoint, error) {
	var points []domain.AvailabilityPoint
	// Iterate day by day
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dayEnd := d.AddDate(0, 0, 1)
		avail, err := r.GetAvailableQuantity(ctx, itemTypeID, d, dayEnd)
		if err != nil {
			return nil, err
		}

		var total int
		err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM assets WHERE item_type_id = $1 AND status != 'retired'", itemTypeID).Scan(&total)
		if err != nil {
			return nil, err
		}

		points = append(points, domain.AvailabilityPoint{
			Date:      d,
			Available: avail,
			Total:     total,
		})
	}
	return points, nil
}

// GetShortageAlerts identifies future bottlenecks in inventory.
func (r *SqlRepository) GetShortageAlerts(ctx context.Context) ([]domain.ShortageAlert, error) {
	// Simplified: Check next 14 days for all item types
	var alerts []domain.ShortageAlert
	start := time.Now()
	end := start.AddDate(0, 0, 14)

	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM item_types WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itID int64
		var itName string
		if err := rows.Scan(&itID, &itName); err != nil {
			continue
		}

		timeline, err := r.GetAvailabilityTimeline(ctx, itID, start, end)
		if err != nil {
			continue
		}

		for _, p := range timeline {
			if p.Available < 0 {
				alerts = append(alerts, domain.ShortageAlert{
					ItemTypeID:    itID,
					ItemTypeName:  itName,
					Date:          p.Date,
					ShortageCount: -p.Available,
					TotalNeeded:   p.Total - p.Available,
					TotalOwned:    p.Total,
				})
			}
		}
	}

	return alerts, nil
}

// GetMaintenanceForecast predicts inspection needs based on calendar cycles and usage.
func (r *SqlRepository) GetMaintenanceForecast(ctx context.Context) ([]domain.MaintenanceForecast, error) {
	// Assets not inspected in > 90 days OR nearing usage limit
	query := `SELECT id, asset_tag, last_inspection_at, usage_hours, next_service_hours FROM assets 
	          WHERE status != 'retired'`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forecasts []domain.MaintenanceForecast
	for rows.Next() {
		var f domain.MaintenanceForecast
		var lastIns *time.Time
		var usage, nextService float64
		if err := rows.Scan(&f.AssetID, &f.AssetTag, &lastIns, &usage, &nextService); err != nil {
			continue
		}

		// Calculate urgency by calendar
		urgencyCalendar := 0.0
		if lastIns == nil {
			urgencyCalendar = 1.0
		} else {
			daysSince := time.Since(*lastIns).Hours() / 24
			urgencyCalendar = daysSince / 90.0
		}

		// Calculate urgency by usage
		urgencyUsage := 0.0
		if nextService > 0 {
			urgencyUsage = usage / nextService
		}

		// Combined urgency (weighted)
		f.UrgencyScore = urgencyCalendar
		if urgencyUsage > f.UrgencyScore {
			f.UrgencyScore = urgencyUsage
		}

		if f.UrgencyScore >= 0.8 {
			if urgencyUsage > urgencyCalendar {
				f.Reason = "Usage limit approached"
			} else {
				f.Reason = "Quarterly cycle exceeded"
			}
			f.NextService = time.Now() // Simplified
			forecasts = append(forecasts, f)
		}
	}
	return forecasts, nil
}

// GetDashboardStats returns a summary of the system state.
func (r *SqlRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{
		AssetsByStatus: make(map[string]int),
	}

	// Total Assets and Group by Status
	queryAssets := "SELECT status, COUNT(*) FROM assets GROUP BY status"
	rows, err := r.db.QueryContext(ctx, queryAssets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.AssetsByStatus[status] = count
		stats.TotalAssets += count
	}

	// Active Rentals (Confirmed or Partially Fulfilled reservations)
	queryRentals := "SELECT COUNT(*) FROM rental_reservations WHERE reservation_status IN ('ReservationConfirmed', 'ReservationPartiallyFulfilled')"
	err = r.db.QueryRowContext(ctx, queryRentals).Scan(&stats.ActiveRentals)
	if err != nil {
		return nil, err
	}

	// Pending Outbox
	queryOutbox := "SELECT COUNT(*) FROM outbox_events WHERE status = 'pending'"
	err = r.db.QueryRowContext(ctx, queryOutbox).Scan(&stats.PendingOutbox)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// BulkRecallAssets transitions multiple assets to 'recalled' status.
func (r *SqlRepository) BulkRecallAssets(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query := "UPDATE assets SET status = 'recalled', updated_at = $1 WHERE id = ANY($2)"
	_, err := r.db.ExecContext(ctx, query, time.Now(), pq.Array(ids))
	if err != nil {
		return fmt.Errorf("bulk recall assets: %w", err)
	}
	return nil
}

// Companies

func (r *SqlRepository) CreateCompany(ctx context.Context, c *domain.Company) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	query := `INSERT INTO companies (name, legal_name, description, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, query, c.Name, c.LegalName, c.Description, c.Metadata, c.CreatedAt, c.UpdatedAt).Scan(&c.ID)
}

func (r *SqlRepository) GetCompany(ctx context.Context, id int64) (*domain.Company, error) {
	query := `SELECT id, name, legal_name, description, metadata, created_at, updated_at FROM companies WHERE id = $1`
	var c domain.Company
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.LegalName, &c.Description, &c.Metadata, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (r *SqlRepository) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, legal_name, description, metadata, created_at, updated_at FROM companies ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.Company
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.LegalName, &c.Description, &c.Metadata, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}

func (r *SqlRepository) UpdateCompany(ctx context.Context, c *domain.Company) error {
	c.UpdatedAt = time.Now()
	query := `UPDATE companies SET name = $1, legal_name = $2, description = $3, metadata = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.LegalName, c.Description, c.Metadata, c.UpdatedAt, c.ID)
	return err
}

// Contacts

// People & Roles

func (r *SqlRepository) CreatePerson(ctx context.Context, p *domain.Person) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	query := `INSERT INTO people (given_name, family_name, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, p.GivenName, p.FamilyName, p.Metadata, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
	if err != nil {
		return err
	}
	for i := range p.ContactPoints {
		cp := &p.ContactPoints[i]
		_, err := r.db.ExecContext(ctx, `INSERT INTO contact_points (person_id, email, phone, contact_type) VALUES ($1, $2, $3, $4)`, p.ID, cp.Email, cp.Phone, cp.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SqlRepository) GetPerson(ctx context.Context, id int64) (*domain.Person, error) {
	var p domain.Person
	err := r.db.QueryRowContext(ctx, `SELECT id, given_name, family_name, metadata, created_at, updated_at FROM people WHERE id = $1`, id).Scan(&p.ID, &p.GivenName, &p.FamilyName, &p.Metadata, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Fetch matching contact points
	rows, err := r.db.QueryContext(ctx, `SELECT email, phone, contact_type FROM contact_points WHERE person_id = $1`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cp domain.ContactPoint
			if err := rows.Scan(&cp.Email, &cp.Phone, &cp.Type); err == nil {
				p.ContactPoints = append(p.ContactPoints, cp)
			}
		}
	}
	// Fetch primary role (latest one)
	_ = r.db.QueryRowContext(ctx, `SELECT organization_id, role_name FROM organization_roles WHERE person_id = $1 ORDER BY id DESC LIMIT 1`, id).Scan(&p.CompanyID, &p.RoleName)

	return &p, nil
}

func (r *SqlRepository) ListPeople(ctx context.Context) ([]domain.Person, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, given_name, family_name, metadata, created_at, updated_at FROM people ORDER BY family_name, given_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Person
	personMap := make(map[int64]*domain.Person)

	for rows.Next() {
		var p domain.Person
		if err := rows.Scan(&p.ID, &p.GivenName, &p.FamilyName, &p.Metadata, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.ContactPoints = []domain.ContactPoint{} // Initialize empty slice
		results = append(results, p)
	}

	for i := range results {
		personMap[results[i].ID] = &results[i]
	}

	// Fetch all contact points and map them
	cpRows, err := r.db.QueryContext(ctx, `SELECT person_id, email, phone, contact_type FROM contact_points`)
	if err == nil {
		defer cpRows.Close()
		for cpRows.Next() {
			var personID int64
			var cp domain.ContactPoint
			if err := cpRows.Scan(&personID, &cp.Email, &cp.Phone, &cp.Type); err == nil {
				if p, ok := personMap[personID]; ok {
					p.ContactPoints = append(p.ContactPoints, cp)
				}
			}
		}
	}

	// Fetch all roles and map them (latest one per person)
	roleRows, err := r.db.QueryContext(ctx, `SELECT person_id, organization_id, role_name FROM organization_roles`)
	if err == nil {
		defer roleRows.Close()
		for roleRows.Next() {
			var personID int64
			var orgID int64
			var roleName string
			if err := roleRows.Scan(&personID, &orgID, &roleName); err == nil {
				if p, ok := personMap[personID]; ok {
					// For now, just take one role as primary
					p.CompanyID = &orgID
					p.RoleName = roleName
				}
			}
		}
	}

	return results, nil
}

func (r *SqlRepository) UpdatePerson(ctx context.Context, p *domain.Person) error {
	p.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE people SET given_name = $1, family_name = $2, metadata = $3, updated_at = $4 WHERE id = $5`, p.GivenName, p.FamilyName, p.Metadata, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	// Sync contact points
	_, _ = r.db.ExecContext(ctx, `DELETE FROM contact_points WHERE person_id = $1`, p.ID)
	for _, cp := range p.ContactPoints {
		_, _ = r.db.ExecContext(ctx, `INSERT INTO contact_points (person_id, email, phone, contact_type) VALUES ($1, $2, $3, $4)`, p.ID, cp.Email, cp.Phone, cp.Type)
	}
	// Sync role if provided
	if p.CompanyID != nil && p.RoleName != "" {
		// Delete old roles and insert primary one (simplification)
		_, _ = r.db.ExecContext(ctx, `DELETE FROM organization_roles WHERE person_id = $1`, p.ID)
		_, _ = r.db.ExecContext(ctx, `INSERT INTO organization_roles (person_id, organization_id, role_name) VALUES ($1, $2, $3)`, p.ID, *p.CompanyID, p.RoleName)
	}
	return nil
}

func (r *SqlRepository) DeletePerson(ctx context.Context, id int64) error {
	_, _ = r.db.ExecContext(ctx, `DELETE FROM contact_points WHERE person_id = $1`, id)
	_, err := r.db.ExecContext(ctx, "DELETE FROM people WHERE id = $1", id)
	return err
}

func (r *SqlRepository) CreateOrganizationRole(ctx context.Context, or *domain.OrganizationRole) error {
	query := `INSERT INTO organization_roles (person_id, organization_id, role_name, start_date, end_date, metadata)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, query, or.PersonID, or.OrganizationID, or.RoleName, or.StartDate, or.EndDate, or.Metadata).Scan(&or.ID)
}

func (r *SqlRepository) ListOrganizationRoles(ctx context.Context, orgID *int64, personID *int64) ([]domain.OrganizationRole, error) {
	query := `SELECT id, person_id, organization_id, role_name, start_date, end_date, metadata FROM organization_roles WHERE 1=1`
	var args []interface{}
	idx := 1
	if orgID != nil {
		query += fmt.Sprintf(` AND organization_id = $%d`, idx)
		args = append(args, *orgID)
		idx++
	}
	if personID != nil {
		query += fmt.Sprintf(` AND person_id = $%d`, idx)
		args = append(args, *personID)
		idx++
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.OrganizationRole
	for rows.Next() {
		var or domain.OrganizationRole
		if err := rows.Scan(&or.ID, &or.PersonID, &or.OrganizationID, &or.RoleName, &or.StartDate, &or.EndDate, &or.Metadata); err != nil {
			return nil, err
		}
		results = append(results, or)
	}
	return results, nil
}

func (r *SqlRepository) DeleteOrganizationRole(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM organization_roles WHERE id = $1", id)
	return err
}

// Places

func (r *SqlRepository) CreatePlace(ctx context.Context, p *domain.Place) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	addrJSON, _ := json.Marshal(p.Address)
	query := `INSERT INTO places (name, description, contained_in_place_id, owner_id, category, address, is_internal, presumed_demands, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.ContainedInPlaceID, p.OwnerID, p.Category, addrJSON, p.IsInternal, p.PresumedDemands, p.Metadata, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
}

func (r *SqlRepository) GetPlace(ctx context.Context, id int64) (*domain.Place, error) {
	query := `SELECT id, name, description, contained_in_place_id, owner_id, category, address, is_internal, presumed_demands, metadata, created_at, updated_at FROM places WHERE id = $1`
	var p domain.Place
	var addrJSON, demandsJSON, metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.ContainedInPlaceID, &p.OwnerID, &p.Category, &addrJSON, &p.IsInternal, &demandsJSON, &metadataJSON, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(addrJSON) > 0 {
		json.Unmarshal(addrJSON, &p.Address)
	}
	p.PresumedDemands = json.RawMessage(demandsJSON)
	p.Metadata = json.RawMessage(metadataJSON)
	return &p, nil
}

func (r *SqlRepository) ListPlaces(ctx context.Context, ownerID *int64, parentID *int64) ([]domain.Place, error) {
	query := `SELECT id, name, description, contained_in_place_id, owner_id, category, address, is_internal, presumed_demands, metadata, created_at, updated_at FROM places WHERE 1=1`
	var args []interface{}
	idx := 1
	if ownerID != nil {
		query += fmt.Sprintf(` AND owner_id = $%d`, idx)
		args = append(args, *ownerID)
		idx++
	}
	if parentID != nil {
		query += fmt.Sprintf(` AND contained_in_place_id = $%d`, idx)
		args = append(args, *parentID)
		idx++
	}
	query += ` ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.Place
	for rows.Next() {
		var p domain.Place
		var addrJSON, demandsJSON, metadataJSON []byte
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.ContainedInPlaceID, &p.OwnerID, &p.Category, &addrJSON, &p.IsInternal, &demandsJSON, &metadataJSON, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if len(addrJSON) > 0 {
			json.Unmarshal(addrJSON, &p.Address)
		}
		p.PresumedDemands = json.RawMessage(demandsJSON)
		p.Metadata = json.RawMessage(metadataJSON)
		results = append(results, p)
	}
	return results, nil
}

func (r *SqlRepository) UpdatePlace(ctx context.Context, p *domain.Place) error {
	p.UpdatedAt = time.Now()
	addrJSON, _ := json.Marshal(p.Address)
	query := `UPDATE places SET name = $1, description = $2, contained_in_place_id = $3, owner_id = $4, category = $5, address = $6, is_internal = $7, presumed_demands = $8, metadata = $9, updated_at = $10 WHERE id = $11`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.ContainedInPlaceID, p.OwnerID, p.Category, addrJSON, p.IsInternal, p.PresumedDemands, p.Metadata, p.UpdatedAt, p.ID)
	return err
}

func (r *SqlRepository) DeletePlace(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM places WHERE id = $1", id)
	return err
}

// Events

func (r *SqlRepository) CreateEvent(ctx context.Context, e *domain.Event) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	query := `INSERT INTO events (company_id, name, description, start_time, end_time, status, parent_event_id, recurrence_rule, last_confirmed_at, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	return r.db.QueryRowContext(ctx, query, e.CompanyID, e.Name, e.Description, e.StartTime, e.EndTime, e.Status, e.ParentEventID, e.RecurrenceRule, e.LastConfirmedAt, e.Metadata, e.CreatedAt, e.UpdatedAt).Scan(&e.ID)
}

func (r *SqlRepository) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	query := `SELECT id, company_id, name, description, start_time, end_time, status, parent_event_id, recurrence_rule, last_confirmed_at, metadata, created_at, updated_at FROM events WHERE id = $1`
	var e domain.Event
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.CompanyID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.ParentEventID, &e.RecurrenceRule, &e.LastConfirmedAt, &e.Metadata, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &e, err
}

func (r *SqlRepository) ListEvents(ctx context.Context, companyID *int64) ([]domain.Event, error) {
	query := `SELECT id, company_id, name, description, start_time, end_time, status, parent_event_id, recurrence_rule, last_confirmed_at, metadata, created_at, updated_at FROM events`
	var args []interface{}
	if companyID != nil {
		query += ` WHERE company_id = $1`
		args = append(args, *companyID)
	}
	query += ` ORDER BY start_time`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Status, &e.ParentEventID, &e.RecurrenceRule, &e.LastConfirmedAt, &e.Metadata, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (r *SqlRepository) UpdateEvent(ctx context.Context, e *domain.Event) error {
	e.UpdatedAt = time.Now()
	query := `UPDATE events SET name = $1, description = $2, start_time = $3, end_time = $4, status = $5, parent_event_id = $6, recurrence_rule = $7, last_confirmed_at = $8, metadata = $9, updated_at = $10 WHERE id = $11`
	_, err := r.db.ExecContext(ctx, query, e.Name, e.Description, e.StartTime, e.EndTime, e.Status, e.ParentEventID, e.RecurrenceRule, e.LastConfirmedAt, e.Metadata, e.UpdatedAt, e.ID)
	return err
}

// EventAssetNeeds

func (r *SqlRepository) CreateEventAssetNeed(ctx context.Context, ean *domain.EventAssetNeed) error {
	now := time.Now()
	ean.CreatedAt = now
	ean.UpdatedAt = now
	query := `INSERT INTO event_asset_needs (event_id, item_type_id, quantity, is_assumed, place_id, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRowContext(ctx, query, ean.EventID, ean.ItemTypeID, ean.Quantity, ean.IsAssumed, ean.PlaceID, ean.Metadata, ean.CreatedAt, ean.UpdatedAt).Scan(&ean.ID)
}

func (r *SqlRepository) ListEventAssetNeeds(ctx context.Context, eventID int64) ([]domain.EventAssetNeed, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, event_id, item_type_id, quantity, is_assumed, place_id, metadata, created_at, updated_at FROM event_asset_needs WHERE event_id = $1`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []domain.EventAssetNeed
	for rows.Next() {
		var ean domain.EventAssetNeed
		if err := rows.Scan(&ean.ID, &ean.EventID, &ean.ItemTypeID, &ean.Quantity, &ean.IsAssumed, &ean.PlaceID, &ean.Metadata, &ean.CreatedAt, &ean.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, ean)
	}
	return results, nil
}

func (r *SqlRepository) UpdateEventAssetNeed(ctx context.Context, ean *domain.EventAssetNeed) error {
	ean.UpdatedAt = time.Now()
	query := `UPDATE event_asset_needs SET item_type_id = $1, quantity = $2, is_assumed = $3, place_id = $4, metadata = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query, ean.ItemTypeID, ean.Quantity, ean.IsAssumed, ean.PlaceID, ean.Metadata, ean.UpdatedAt, ean.ID)
	return err
}

// Delete methods for entities

func (r *SqlRepository) DeleteCompany(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM companies WHERE id = $1", id)
	return err
}

func (r *SqlRepository) DeleteEvent(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	return err
}

// GetDefaultInternalPlace returns the system default internal warehouse or creates one if it doesn't exist.
func (r *SqlRepository) GetDefaultInternalPlace(ctx context.Context) (*domain.Place, error) {
	// 1. Try to find an existing internal place named 'Main Warehouse'
	var p domain.Place
	var addrJSON, demandsJSON, metadataJSON []byte
	query := `SELECT id, name, description, contained_in_place_id, owner_id, category, address, is_internal, presumed_demands, metadata, created_at, updated_at 
	          FROM places WHERE is_internal = TRUE LIMIT 1`
	err := r.db.QueryRowContext(ctx, query).Scan(
		&p.ID, &p.Name, &p.Description, &p.ContainedInPlaceID, &p.OwnerID, &p.Category, &addrJSON, &p.IsInternal, &demandsJSON, &metadataJSON, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == nil {
		if len(addrJSON) > 0 {
			json.Unmarshal(addrJSON, &p.Address)
		}
		p.PresumedDemands = json.RawMessage(demandsJSON)
		p.Metadata = json.RawMessage(metadataJSON)
		return &p, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// 2. Not found, create a new default place
	p = domain.Place{
		Name:       "Main Warehouse",
		IsInternal: true,
		Category:   func(s string) *string { return &s }("site"),
	}
	err = r.CreatePlace(ctx, &p)
	return &p, err
}
