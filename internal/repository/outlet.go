package repository

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/response"
)

type OutletRepository struct {
	db *sql.DB
}

func NewOutletRepository(db *sql.DB) *OutletRepository {
	return &OutletRepository{
		db: db,
	}
}

func (r *OutletRepository) Create(ctx context.Context, tx *sql.Tx, outlet *model.Outlet) (*dto.OutletResponse, error) {
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

	result, err := tx.ExecContext(
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	res := &dto.OutletResponse{
		OutletID: uint(id),
		Name:     outlet.Name,
		Address:  outlet.Address,
		Channel:  outlet.Channel,
		Tier:     outlet.Tier,
		SalesID:  outlet.SalesID,
		Stage:    outlet.Stage,
		Note:     *outlet.Note,
	}

	return res, nil
}

func (r *OutletRepository) FindByID(ctx context.Context, tx *sql.Tx, id uint) (*model.Outlet, error) {
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

	var outlet model.Outlet

	err := tx.QueryRowContext(
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

func (r *OutletRepository) FindAll(ctx context.Context, tx *sql.Tx, pagination *response.Pagination) ([]model.Outlet, error) {
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM outlets
	`
	err := tx.QueryRowContext(
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

	rows, err := tx.QueryContext(ctx, query, pagination.PerPage, (pagination.Page-1)*pagination.PerPage)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	outlets := make([]model.Outlet, 0)

	for rows.Next() {
		var outlet model.Outlet

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

func (r *OutletRepository) Update(ctx context.Context, tx *sql.Tx, outlet *model.Outlet) error {
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

	_, err := tx.ExecContext(
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

func (r *OutletRepository) Delete(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `
		DELETE FROM outlets
		WHERE outlet_id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}
