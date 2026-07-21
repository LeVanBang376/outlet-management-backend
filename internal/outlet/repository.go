package outlet

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/response"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, outlet *Outlet) error {
	query := `
		INSERT INTO outlets (
			name,
			address,
			channel,
			tier,
			sales_id,
			stage,
			note
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		outlet.Name,
		outlet.Address,
		outlet.Channel,
		outlet.Tier,
		outlet.SalesID,
		outlet.Stage,
		outlet.Note,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	outlet.OutletID = uint(id)

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uint) (*Outlet, error) {
	query := `
		SELECT 
			outlet_id,
			name,
			address,
			channel,
			tier,
			sales_id,
			stage,
			note
		FROM outlets
		WHERE outlet_id = ?
	`

	var outlet Outlet

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&outlet.OutletID,
		&outlet.Name,
		&outlet.Address,
		&outlet.Channel,
		&outlet.Tier,
		&outlet.SalesID,
		&outlet.Stage,
		&outlet.Note,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &outlet, nil
}

func (r *Repository) FindAll(ctx context.Context, pagination *response.Pagination) ([]Outlet, error) {
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM outlets
	`
	err := r.db.QueryRowContext(
		ctx,
		countQuery,
	).Scan(&total)
	if err != nil {
		return nil, err
	}

	pagination.SetTotal(total)

	query := `
		SELECT
			outlet_id,
			name,
			address,
			channel,
			tier,
			sales_id,
			stage,
			note
		FROM outlets
		ORDER BY outlet_id DESC
		LIMIT ?
		OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, pagination.PerPage, (pagination.Page-1)*pagination.PerPage)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	outlets := make([]Outlet, 0)

	for rows.Next() {
		var outlet Outlet

		err := rows.Scan(
			&outlet.OutletID,
			&outlet.Name,
			&outlet.Address,
			&outlet.Channel,
			&outlet.Tier,
			&outlet.SalesID,
			&outlet.Stage,
			&outlet.Note,
		)

		if err != nil {
			return nil, err
		}

		outlets = append(outlets, outlet)
	}

	return outlets, rows.Err()
}

func (r *Repository) Update(ctx context.Context, outlet *Outlet) error {
	query := `
		UPDATE outlets
		SET
			name = ?,
			address = ?,
			channel = ?,
			tier = ?,
			sales_id = ?,
			stage = ?,
			note = ?
		WHERE outlet_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		outlet.Name,
		outlet.Address,
		outlet.Channel,
		outlet.Tier,
		outlet.SalesID,
		outlet.Stage,
		outlet.Note,
		outlet.OutletID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	query := `
		DELETE FROM outlets
		WHERE outlet_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}
