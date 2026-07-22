package repository

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/model"
)

type SalesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{
		db: db,
	}
}

func (r *SalesRepository) GetAll(ctx context.Context, tx *sql.Tx) ([]model.Sales, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT 
			sales_id,
			name
		FROM sales
		ORDER BY sales_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Sales, 0)

	for rows.Next() {
		var item model.Sales

		if err := rows.Scan(
			&item.SalesID,
			&item.Name,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}
