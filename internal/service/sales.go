package service

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/repository"
)

type SalesService struct {
	db   *sql.DB
	repo *repository.SalesRepository
}

func NewSalesService(db *sql.DB, repository *repository.SalesRepository) *SalesService {
	return &SalesService{
		db:   db,
		repo: repository,
	}
}

func (s *SalesService) GetAll(ctx context.Context) ([]dto.SalesResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sales, err := s.repo.GetAll(ctx, tx)

	if err != nil {
		return nil, err
	}

	result := make([]dto.SalesResponse, 0, len(sales))

	for _, item := range sales {
		result = append(result, dto.NewSalesResponse(item))
	}

	return result, tx.Commit()
}
