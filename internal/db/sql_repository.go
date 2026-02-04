package db

import (
	"context"
	"database/sql"
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

	query := `INSERT INTO item_types (code, name, kind, is_active, schema_org, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, it.SchemaOrg, it.Metadata, it.CreatedAt, it.UpdatedAt,
	).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("create item_type: %w", err)
	}
	return nil
}

// GetItemTypeByID retrieves an item type by its ID.
func (r *SqlRepository) GetItemTypeByID(ctx context.Context, id int64) (*domain.ItemType, error) {
	query := `SELECT id, code, name, kind, is_active, schema_org, metadata, created_at, updated_at 
	          FROM item_types WHERE id = $1`

	var it domain.ItemType
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &it.SchemaOrg, &it.Metadata, &it.CreatedAt, &it.UpdatedAt,
	)
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
		if err := rows.Scan(&it.ID, &it.Code, &it.Name, &it.Kind, &it.IsActive, &it.SchemaOrg, &it.Metadata, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item_type: %w", err)
		}
		results = append(results, it)
	}
	return results, nil
}

// UpdateItemType updates an existing item type.
func (r *SqlRepository) UpdateItemType(ctx context.Context, it *domain.ItemType) error {
	it.UpdatedAt = time.Now()
	query := `UPDATE item_types SET code = $1, name = $2, kind = $3, is_active = $4, schema_org = $5, metadata = $6, updated_at = $7
	          WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		it.Code, it.Name, it.Kind, it.IsActive, it.SchemaOrg, it.Metadata, it.UpdatedAt, it.ID,
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
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE id = $1`

	var a domain.Asset
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.SchemaOrg, &a.Metadata, &a.CreatedAt, &a.UpdatedAt,
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

	query := `INSERT INTO assets (item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, schema_org, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo, a.MeshNodeID, a.WireguardHostname, a.SchemaOrg, a.Metadata, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.ID)
	if err != nil {
		return fmt.Errorf("create asset: %w", err)
	}
	return nil
}

// ListAssetsByItemType returns assets belonging to a specific item type.
func (r *SqlRepository) ListAssetsByItemType(ctx context.Context, itemTypeID int64) ([]domain.Asset, error) {
	query := `SELECT id, item_type_id, asset_tag, serial_number, status, location, assigned_to, mesh_node_id, wireguard_hostname, schema_org, metadata, created_at, updated_at 
	          FROM assets WHERE item_type_id = $1`

	rows, err := r.db.QueryContext(ctx, query, itemTypeID)
	if err != nil {
		return nil, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	var results []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.ItemTypeID, &a.AssetTag, &a.SerialNumber, &a.Status, &a.Location, &a.AssignedTo, &a.MeshNodeID, &a.WireguardHostname, &a.SchemaOrg, &a.Metadata, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		results = append(results, a)
	}
	return results, nil
}

// UpdateAsset updates an existing asset.
func (r *SqlRepository) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	a.UpdatedAt = time.Now()
	query := `UPDATE assets SET item_type_id = $1, asset_tag = $2, serial_number = $3, status = $4, location = $5, assigned_to = $6, mesh_node_id = $7, wireguard_hostname = $8, schema_org = $9, metadata = $10, updated_at = $11
	          WHERE id = $12`

	_, err := r.db.ExecContext(ctx, query,
		a.ItemTypeID, a.AssetTag, a.SerialNumber, a.Status, a.Location, a.AssignedTo, a.MeshNodeID, a.WireguardHostname, a.SchemaOrg, a.Metadata, a.UpdatedAt, a.ID,
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
