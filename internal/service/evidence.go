package service

import (
	"context"
	"database/sql"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
	"time"
)

type EvidenceService struct {
	db   *sql.DB
	repo *repository.EvidenceRepository
}

func NewEvidenceService(db *sql.DB, repo *repository.EvidenceRepository) *EvidenceService {
	return &EvidenceService{
		db:   db,
		repo: repo,
	}
}

func (s *EvidenceService) Create(
	ctx context.Context,
	req dto.CreateEvidenceRequest,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	evidence := model.Evidence{
		ScheduleID:  req.ScheduleID,
		ObjectKey:   req.ObjectKey,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Size:        req.Size,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(
		ctx,
		tx,
		&evidence,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *EvidenceService) GetByID(
	ctx context.Context,
	id uint,
) (*model.Evidence, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	evidence, err := s.repo.FindByID(
		ctx,
		tx,
		id,
	)

	if err != nil {
		return nil, err
	}

	if evidence == nil {
		return nil, customerrors.EvidenceErrNotFound
	}

	return evidence, tx.Commit()
}

func (s *EvidenceService) GetByScheduleID(
	ctx context.Context,
	scheduleID uint,
) ([]model.Evidence, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return []model.Evidence{}, err
	}
	defer tx.Rollback()

	evidences, err := s.repo.FindByScheduleID(
		ctx,
		tx,
		scheduleID,
	)
	if err != nil {
		return nil, err
	}
	return evidences, tx.Commit()
}

func (s *EvidenceService) Delete(
	ctx context.Context,
	id uint,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	evidence, err := s.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if evidence == nil {
		return customerrors.EvidenceErrNotFound
	}

	if err := s.repo.Delete(
		ctx,
		tx,
		id,
	); err != nil {
		return err
	}

	return tx.Commit()
}
