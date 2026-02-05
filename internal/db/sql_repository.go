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
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, management_url,
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 usage_hours, next_service_hours, created_by_user_id, updated_by_user_id, schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE id = $1`

	var a domain.Asset
	var schemaOrgJSON, metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
		&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
		&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == nil {
		a.SchemaOrg = json.RawMessage(schemaOrgJSON)
		a.Metadata = json.RawMessage(metadataJSON)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan asset: %w", err)
	}
	return &a, nil
}

// CreateAsset creates a new asset.
func (r *SqlRepository) CreateAsset(ctx context.Context, a *domain.Asset) error {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	query := `INSERT INTO assets (
		item_type_id, asset_tag, serial_number, status, location, assigned_to, 
		mesh_node_id, wireguard_hostname, management_url, build_spec_version, provisioning_status, 
		firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
		usage_hours, next_service_hours, schema_org, metadata, created_by_user_id, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo,
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
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, management_url, 
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
			&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
			&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
			&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		a.SchemaOrg = json.RawMessage(schemaOrgJSON)
		a.Metadata = json.RawMessage(metadataJSON)
		results = append(results, a)
	}
	return results, nil
}

// ListAssetsByItemType returns assets belonging to a specific item type.
func (r *SqlRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, management_url, 
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
			&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.ManagementURL,
			&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
			&a.UsageHours, &a.NextServiceHours, &a.CreatedByUserID, &a.UpdatedByUserID, &schemaOrgJSON, &metadataJSON, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		a.SchemaOrg = json.RawMessage(schemaOrgJSON)
		a.Metadata = json.RawMessage(metadataJSON)
		results = append(results, a)
	}
	return results, nil
}

// UpdateAsset updates an existing asset.
func (r *SqlRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	a.UpdatedAt = time.Now()
	query := `UPDATE assets SET 
		item_type_id = $1, asset_tag = $2, serial_number = $3, status = $4, 
		location = $5, assigned_to = $6, mesh_node_id = $7, wireguard_hostname = $8,
		management_url = $9, build_spec_version = $10, provisioning_status = $11, firmware_version = $12,
		hostname = $13, remote_management_id = $14, current_build_spec_id = $15, last_inspection_at = $16,
		usage_hours = $17, next_service_hours = $18, updated_by_user_id = $19, schema_org = $20, 
		metadata = $21, updated_at = $22
		WHERE id = $23`

	_, err := r.db.ExecContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo,
		a.MeshNodeID, a.WireguardHostname, a.ManagementURL, a.BuildSpecVersion, a.ProvisioningStatus,
		a.FirmwareVersion, a.Hostname, a.RemoteManagementID, a.CurrentBuildSpecID, a.LastInspectionAt,
		a.UsageHours, a.NextServiceHours, a.UpdatedByUserID, a.SchemaOrg, a.Metadata, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return fmt.Errorf("update asset: %w", err)
	}
	return nil
}

// UpdateAssetStatus updates the status of an asset.
func (r *SqlRepository) UpdateAssetStatus(ctx context.Context, id int64, status domain.AssetStatus) error {
	query := `UPDATE assets SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
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

// CreateRentAction creates a new rent action and its associated items.
func (r *SqlRepository) CreateRentAction(ctx context.Context, ra *domain.RentAction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	ra.CreatedAt = now
	ra.UpdatedAt = now

	query := `INSERT INTO rent_actions (
		requester_ref, created_by_ref, approved_by_ref, status, priority, 
		start_time, end_time, is_asap, description, external_source, 
		external_ref, schema_org, metadata, created_by_user_id, updated_by_user_id, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		ra.RequesterRef, ra.CreatedByRef, ra.ApprovedByRef, ra.Status, ra.Priority,
		ra.StartTime, ra.EndTime, ra.IsASAP, ra.Description, ra.ExternalSource,
		ra.ExternalRef, ra.SchemaOrg, ra.Metadata, ra.CreatedByUserID, ra.UpdatedByUserID, ra.CreatedAt, ra.UpdatedAt,
	).Scan(&ra.ID)
	if err != nil {
		return fmt.Errorf("insert rent_action: %w", err)
	}

	for i := range ra.Items {
		item := &ra.Items[i]
		item.RentActionID = ra.ID
		itemQuery := `INSERT INTO rent_action_items (
			rent_action_id, item_kind, item_id, requested_quantity, allocated_quantity, notes, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

		err = tx.QueryRowContext(ctx, itemQuery,
			item.RentActionID, item.ItemKind, item.ItemID, item.RequestedQuantity,
			item.AllocatedQuantity, item.Notes, item.Metadata,
		).Scan(&item.ID)
		if err != nil {
			return fmt.Errorf("insert rent_action_item: %w", err)
		}
	}

	return tx.Commit()
}

// GetRentActionByID retrieves a rent action by its ID, including its line items.
func (r *SqlRepository) GetRentActionByID(ctx context.Context, id int64) (*domain.RentAction, error) {
	query := `SELECT id, requester_ref, created_by_ref, approved_by_ref, status, priority, 
	                 start_time, end_time, is_asap, description, external_source, 
	                 external_ref, schema_org, metadata, approved_at, rejected_at, 
	                 cancelled_at, created_at, updated_at 
	          FROM rent_actions WHERE id = $1`

	var ra domain.RentAction
	var schemaOrgJSON, metadataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ra.ID, &ra.RequesterRef, &ra.CreatedByRef, &ra.ApprovedByRef, &ra.Status, &ra.Priority,
		&ra.StartTime, &ra.EndTime, &ra.IsASAP, &ra.Description, &ra.ExternalSource,
		&ra.ExternalRef, &schemaOrgJSON, &metadataJSON, &ra.ApprovedAt, &ra.RejectedAt,
		&ra.CancelledAt, &ra.CreatedAt, &ra.UpdatedAt,
	)
	if err == nil {
		ra.SchemaOrg = json.RawMessage(schemaOrgJSON)
		ra.Metadata = json.RawMessage(metadataJSON)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan rent_action: %w", err)
	}

	itemQuery := `SELECT id, rent_action_id, item_kind, item_id, requested_quantity, allocated_quantity, notes, metadata 
	              FROM rent_action_items WHERE rent_action_id = $1`

	rows, err := r.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return nil, fmt.Errorf("query rent_action_items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.RentActionItem
		var metadataJSON []byte
		if err := rows.Scan(&item.ID, &item.RentActionID, &item.ItemKind, &item.ItemID, &item.RequestedQuantity, &item.AllocatedQuantity, &item.Notes, &metadataJSON); err != nil {
			return nil, fmt.Errorf("scan rent_action_item: %w", err)
		}
		item.Metadata = json.RawMessage(metadataJSON)
		ra.Items = append(ra.Items, item)
	}

	return &ra, nil
}

// ListRentActions returns all rent actions.
func (r *SqlRepository) ListRentActions(ctx context.Context) ([]domain.RentAction, error) {
	query := `SELECT id, requester_ref, created_by_ref, created_by_user_id, approved_by_ref, status, priority, 
	                 start_time, end_time, is_asap, description, external_source, 
	                 external_ref, schema_org, metadata, approved_at, rejected_at, 
	                 cancelled_at, created_at, updated_at 
	          FROM rent_actions ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rent_actions: %w", err)
	}
	defer rows.Close()

	var results []domain.RentAction
	for rows.Next() {
		var ra domain.RentAction
		var schemaOrgJSON, metadataJSON []byte
		err := rows.Scan(
			&ra.ID, &ra.RequesterRef, &ra.CreatedByRef, &ra.CreatedByUserID, &ra.ApprovedByRef, &ra.Status, &ra.Priority,
			&ra.StartTime, &ra.EndTime, &ra.IsASAP, &ra.Description, &ra.ExternalSource,
			&ra.ExternalRef, &schemaOrgJSON, &metadataJSON, &ra.ApprovedAt, &ra.RejectedAt,
			&ra.CancelledAt, &ra.CreatedAt, &ra.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rent_action: %w", err)
		}
		ra.SchemaOrg = json.RawMessage(schemaOrgJSON)
		ra.Metadata = json.RawMessage(metadataJSON)
		results = append(results, ra)
	}
	return results, nil
}

// UpdateRentAction updates an existing rent action.
func (r *SqlRepository) UpdateRentAction(ctx context.Context, ra *domain.RentAction) error {
	ra.UpdatedAt = time.Now()
	query := `UPDATE rent_actions SET 
		requester_ref = $1, created_by_ref = $2, approved_by_ref = $3, status = $4, 
		priority = $5, start_time = $6, end_time = $7, is_asap = $8, 
		description = $9, external_source = $10, external_ref = $11, 
		schema_org = $12, metadata = $13, approved_at = $14, rejected_at = $15, 
		cancelled_at = $16, updated_at = $17 
		WHERE id = $18`

	_, err := r.db.ExecContext(ctx, query,
		ra.RequesterRef, ra.CreatedByRef, ra.ApprovedByRef, ra.Status,
		ra.Priority, ra.StartTime, ra.EndTime, ra.IsASAP,
		ra.Description, ra.ExternalSource, ra.ExternalRef,
		ra.SchemaOrg, ra.Metadata, ra.ApprovedAt, ra.RejectedAt,
		ra.CancelledAt, ra.UpdatedAt, ra.ID,
	)
	if err != nil {
		return fmt.Errorf("update rent_action: %w", err)
	}
	return nil
}

// UpdateRentActionStatus updates only the status and a specific timestamp field for a rent action.
func (r *SqlRepository) UpdateRentActionStatus(ctx context.Context, id int64, status domain.RentActionStatus, timestampField string, timestampValue time.Time) error {
	var query string
	if timestampField != "" {
		query = fmt.Sprintf("UPDATE rent_actions SET status = $1, %s = $2, updated_at = $2 WHERE id = $3", timestampField)
		_, err := r.db.ExecContext(ctx, query, status, timestampValue, id)
		return err
	}
	query = "UPDATE rent_actions SET status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// GetAvailableQuantity calculates the available inventory for an item type in a given time window.
func (r *SqlRepository) GetAvailableQuantity(ctx context.Context, itemTypeID int64, startTime, endTime time.Time) (int, error) {
	// 1. Get total assets for this item type (excluding retired)
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM assets WHERE item_type_id = $1 AND status != 'retired'", itemTypeID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count assets: %w", err)
	}

	// 2. Subtract quantities from overlapping APPROVED reservations
	// Overlap rule: (startA < endB) AND (endA > startB)
	query := `
		SELECT COALESCE(SUM(rai.requested_quantity), 0)
		FROM rent_action_items rai
		JOIN rent_actions ra ON rai.rent_action_id = ra.id
		WHERE rai.item_kind = 'item_type' 
		  AND rai.item_id = $1
		  AND ra.status = 'approved'
		  AND ra.start_time < $3
		  AND ra.end_time > $2
	`
	var reserved int
	err = r.db.QueryRowContext(ctx, query, itemTypeID, startTime, endTime).Scan(&reserved)
	if err != nil {
		return 0, fmt.Errorf("sum reserved quantity: %w", err)
	}

	return total - reserved, nil
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

// SubmitInspection records a completed inspection form.
func (r *SqlRepository) SubmitInspection(ctx context.Context, is *domain.InspectionSubmission) error {
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

	// Active Rentals
	queryRentals := "SELECT COUNT(*) FROM rent_actions WHERE status = 'approved' OR status = 'ongoing'"
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
