package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
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

	query := `INSERT INTO item_types (code, name, kind, is_active, supported_features, schema_org, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	featuresJSON, _ := json.Marshal(it.SupportedFeatures)
	err := r.db.QueryRowContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, featuresJSON, it.SchemaOrg, it.Metadata, it.CreatedAt, it.UpdatedAt,
	).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("create item_type: %w", err)
	}
	return nil
}

// GetItemTypeByID retrieves an item type by its ID.
func (r *SqlRepository) GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error) {
	query := `SELECT id, code, name, kind, is_active, supported_features, schema_org, metadata, created_at, updated_at 
	          FROM item_types WHERE id = $1`

	var it domain.ItemType
	var featuresJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &featuresJSON, &it.SchemaOrg, &it.Metadata, &it.CreatedAt, &it.UpdatedAt,
	)
	if err == nil {
		json.Unmarshal(featuresJSON, &it.SupportedFeatures)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan item_type: %w", err)
	}
	return &it, nil
}

// ListItemTypes returns all active item types.
func (r *SqlRepository) ListItemTypes(ctx context.Context) ([]domain.ItemType, error) {
	query := `SELECT id, code, name, kind, is_active, schema_org, metadata, created_at, updated_at 
	          FROM item_types WHERE is_active = TRUE`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query item_types: %w", err)
	}
	defer rows.Close()

	var results []domain.ItemType
	for rows.Next() {
		var it domain.ItemType
		var featuresJSON []byte
		if err := rows.Scan(&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &featuresJSON, &it.SchemaOrg, &it.Metadata, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item_type: %w", err)
		}
		json.Unmarshal(featuresJSON, &it.SupportedFeatures)
		results = append(results, it)
	}
	return results, nil
}

// UpdateItemType updates an existing item type.
func (r *SqlRepository) UpdateItemType(ctx context.Context, it *domain.ItemType) error {
	it.UpdatedAt = time.Now()
	query := `UPDATE item_types SET code = $1, name = $2, kind = $3, is_active = $4, supported_features = $5, schema_org = $6, metadata = $7, updated_at = $8
	          WHERE id = $9`

	featuresJSON, _ := json.Marshal(it.SupportedFeatures)
	_, err := r.db.ExecContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, featuresJSON, it.SchemaOrg, it.Metadata, it.UpdatedAt, it.ID,
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
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, 
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE id = $1`

	var a domain.Asset
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname,
		&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
		&a.SchemaOrg, &a.Metadata, &a.CreatedAt, &a.UpdatedAt,
	)
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
		mesh_node_id, wireguard_hostname, build_spec_version, provisioning_status, 
		firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
		schema_org, metadata, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo,
		a.MeshNodeID, a.WireguardHostname, a.BuildSpecVersion, a.ProvisioningStatus,
		a.FirmwareVersion, a.Hostname, a.RemoteManagementID, a.CurrentBuildSpecID, a.LastInspectionAt,
		a.SchemaOrg, a.Metadata, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.ID)
	if err != nil {
		return fmt.Errorf("create asset: %w", err)
	}
	return nil
}

// ListAssetsByItemType returns assets belonging to a specific item type.
func (r *SqlRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, 
	                 build_spec_version, provisioning_status, firmware_version, hostname, remote_management_id, current_build_spec_id, last_inspection_at,
	                 schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE item_type_id = $1`

	rows, err := r.db.QueryContext(ctx, query, itemTypeID)
	if err != nil {
		return nil, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	var results []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(
			&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname,
			&a.BuildSpecVersion, &a.ProvisioningStatus, &a.FirmwareVersion, &a.Hostname, &a.RemoteManagementID, &a.CurrentBuildSpecID, &a.LastInspectionAt,
			&a.SchemaOrg, &a.Metadata, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
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
		location = $5, assigned_to = $6, mesh_node_id = $7, wireguard_hostname = $8,
		build_spec_version = $9, provisioning_status = $10, firmware_version = $11,
		hostname = $12, remote_management_id = $13, current_build_spec_id = $14, last_inspection_at = $15,
		schema_org = $16, metadata = $17, updated_at = $18
		WHERE id = $19`

	_, err := r.db.ExecContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo,
		a.MeshNodeID, a.WireguardHostname, a.BuildSpecVersion, a.ProvisioningStatus,
		a.FirmwareVersion, a.Hostname, a.RemoteManagementID, a.CurrentBuildSpecID, a.LastInspectionAt,
		a.SchemaOrg, a.Metadata, a.UpdatedAt, a.ID,
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
		external_ref, schema_org, metadata, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		ra.RequesterRef, ra.CreatedByRef, ra.ApprovedByRef, ra.Status, ra.Priority,
		ra.StartTime, ra.EndTime, ra.IsASAP, ra.Description, ra.ExternalSource,
		ra.ExternalRef, ra.SchemaOrg, ra.Metadata, ra.CreatedAt, ra.UpdatedAt,
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ra.ID, &ra.RequesterRef, &ra.CreatedByRef, &ra.ApprovedByRef, &ra.Status, &ra.Priority,
		&ra.StartTime, &ra.EndTime, &ra.IsASAP, &ra.Description, &ra.ExternalSource,
		&ra.ExternalRef, &ra.SchemaOrg, &ra.Metadata, &ra.ApprovedAt, &ra.RejectedAt,
		&ra.CancelledAt, &ra.CreatedAt, &ra.UpdatedAt,
	)
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
		if err := rows.Scan(&item.ID, &item.RentActionID, &item.ItemKind, &item.ItemID, &item.RequestedQuantity, &item.AllocatedQuantity, &item.Notes, &item.Metadata); err != nil {
			return nil, fmt.Errorf("scan rent_action_item: %w", err)
		}
		ra.Items = append(ra.Items, item)
	}

	return &ra, nil
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bs.ID, &bs.Version, &bs.HardwareConfig, &bs.SoftwareConfig, &bs.FirmwareURL, &bs.Metadata, &bs.CreatedAt, &bs.UpdatedAt,
	)
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
		if err := rows.Scan(&bs.ID, &bs.Version, &bs.HardwareConfig, &bs.SoftwareConfig, &bs.FirmwareURL, &bs.Metadata, &bs.CreatedAt, &bs.UpdatedAt); err != nil {
			return nil, err
		}
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
