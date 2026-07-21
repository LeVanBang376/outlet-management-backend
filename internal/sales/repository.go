package sales

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]Sales, error) {
	rows, err := r.db.QueryContext(ctx, `
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

	result := make([]Sales, 0)

	for rows.Next() {
		var item Sales

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
