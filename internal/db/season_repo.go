package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
)

func (r *SqlRepository) CreateShowCompany(ctx context.Context, sc *domain.ShowCompany) error {
	now := time.Now()
	sc.CreatedAt = now
	sc.UpdatedAt = now
	query := `INSERT INTO show_companies (company_id, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, sc.CompanyID, sc.Metadata, sc.CreatedAt, sc.UpdatedAt).Scan(&sc.ID)
}

func (r *SqlRepository) GetShowCompany(ctx context.Context, id int64) (*domain.ShowCompany, error) {
	query := `SELECT id, company_id, metadata, created_at, updated_at FROM show_companies WHERE id = $1`
	var sc domain.ShowCompany
	var metaJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&sc.ID, &sc.CompanyID, &metaJSON, &sc.CreatedAt, &sc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sc.Metadata = json.RawMessage(metaJSON)
	return &sc, nil
}

func (r *SqlRepository) CreateSeason(ctx context.Context, s *domain.Season) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	query := `INSERT INTO seasons (show_company_id, name, start_date, end_date, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query, s.ShowCompanyID, s.Name, s.StartDate, s.EndDate, s.Metadata, s.CreatedAt, s.UpdatedAt).Scan(&s.ID)
}

func (r *SqlRepository) ListSeasonsForCompany(ctx context.Context, showCompanyID int64) ([]domain.Season, error) {
	query := `SELECT id, show_company_id, name, start_date, end_date, metadata, created_at, updated_at
	          FROM seasons WHERE show_company_id = $1 ORDER BY start_date DESC`
	rows, err := r.db.QueryContext(ctx, query, showCompanyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Season
	for rows.Next() {
		var s domain.Season
		var metaJSON []byte
		if err := rows.Scan(&s.ID, &s.ShowCompanyID, &s.Name, &s.StartDate, &s.EndDate, &metaJSON, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.Metadata = json.RawMessage(metaJSON)
		results = append(results, s)
	}
	return results, nil
}

func (r *SqlRepository) CreateRing(ctx context.Context, ring *domain.Ring) error {
	now := time.Now()
	ring.CreatedAt = now
	ring.UpdatedAt = now
	query := `INSERT INTO rings (show_company_id, name, description, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, query, ring.ShowCompanyID, ring.Name, ring.Description, ring.Metadata, ring.CreatedAt, ring.UpdatedAt).Scan(&ring.ID)
}

func (r *SqlRepository) ListRingsForCompany(ctx context.Context, showCompanyID int64) ([]domain.Ring, error) {
	query := `SELECT id, show_company_id, name, description, metadata, created_at, updated_at
	          FROM rings WHERE show_company_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, showCompanyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Ring
	for rows.Next() {
		var ring domain.Ring
		var metaJSON []byte
		if err := rows.Scan(&ring.ID, &ring.ShowCompanyID, &ring.Name, &ring.Description, &metaJSON, &ring.CreatedAt, &ring.UpdatedAt); err != nil {
			return nil, err
		}
		ring.Metadata = json.RawMessage(metaJSON)
		results = append(results, ring)
	}
	return results, nil
}

func (r *SqlRepository) CreateShow(ctx context.Context, s *domain.Show) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	query := `INSERT INTO shows (season_id, name, start_date, end_date, location_id, metadata, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRowContext(ctx, query, s.SeasonID, s.Name, s.StartDate, s.EndDate, s.LocationID, s.Metadata, s.CreatedAt, s.UpdatedAt).Scan(&s.ID)
}

func (r *SqlRepository) GetShowByID(ctx context.Context, id int64) (*domain.Show, error) {
	query := `SELECT id, season_id, name, start_date, end_date, location_id, metadata, created_at, updated_at
	          FROM shows WHERE id = $1`
	var s domain.Show
	var metaJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.SeasonID, &s.Name, &s.StartDate, &s.EndDate, &s.LocationID, &metaJSON, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.Metadata = json.RawMessage(metaJSON)

	// Fetch Rings
	rings, err := r.GetRingsForShow(ctx, s.ID)
	if err == nil {
		s.Rings = rings
	}

	return &s, nil
}

func (r *SqlRepository) AddRingToShow(ctx context.Context, sr *domain.ShowRing) error {
	query := `INSERT INTO show_rings (show_id, ring_id) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, sr.ShowID, sr.RingID).Scan(&sr.ID)
}

func (r *SqlRepository) GetRingsForShow(ctx context.Context, showID int64) ([]domain.ShowRing, error) {
	query := `SELECT sr.id, sr.show_id, sr.ring_id, r.id, r.show_company_id, r.name, r.description, r.metadata, r.created_at, r.updated_at
	          FROM show_rings sr
	          JOIN rings r ON sr.ring_id = r.id
	          WHERE sr.show_id = $1`
	rows, err := r.db.QueryContext(ctx, query, showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.ShowRing
	for rows.Next() {
		var sr domain.ShowRing
		var ring domain.Ring
		var metaJSON []byte
		if err := rows.Scan(
			&sr.ID, &sr.ShowID, &sr.RingID,
			&ring.ID, &ring.ShowCompanyID, &ring.Name, &ring.Description, &metaJSON, &ring.CreatedAt, &ring.UpdatedAt,
		); err != nil {
			return nil, err
		}
		ring.Metadata = json.RawMessage(metaJSON)
		sr.Ring = &ring

		// Get Loadout Items
		loadoutQuery := `SELECT id, show_ring_id, item_type_id, quantity FROM ring_loadout_items WHERE show_ring_id = $1`
		lr, err := r.db.QueryContext(ctx, loadoutQuery, sr.ID)
		if err == nil {
			for lr.Next() {
				var item domain.RingLoadoutItem
				if err := lr.Scan(&item.ID, &item.ShowRingID, &item.ItemTypeID, &item.Quantity); err == nil {
					sr.LoadoutItems = append(sr.LoadoutItems, item)
				}
			}
			lr.Close()
		}

		results = append(results, sr)
	}
	return results, nil
}

func (r *SqlRepository) SetShowRingLoadout(ctx context.Context, showRingID int64, items []domain.RingLoadoutItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing
	_, err = tx.ExecContext(ctx, `DELETE FROM ring_loadout_items WHERE show_ring_id = $1`, showRingID)
	if err != nil {
		return fmt.Errorf("failed to clear old loadout: %w", err)
	}

	// Insert new
	for i := range items {
		item := &items[i]
		item.ShowRingID = showRingID
		query := `INSERT INTO ring_loadout_items (show_ring_id, item_type_id, quantity) VALUES ($1, $2, $3) RETURNING id`
		if err := tx.QueryRowContext(ctx, query, item.ShowRingID, item.ItemTypeID, item.Quantity).Scan(&item.ID); err != nil {
			return fmt.Errorf("failed to insert loadout item: %w", err)
		}
	}

	return tx.Commit()
}
