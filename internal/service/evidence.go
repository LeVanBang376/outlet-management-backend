package service

import (
	"context"
	"database/sql"

	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
)

type EvidenceService struct {
	db       *sql.DB
	repo     *repository.EvidenceRepository
	fileRepo *repository.FileRepository
}

func NewEvidenceService(
	db *sql.DB,
	repo *repository.EvidenceRepository,
	fileRepo *repository.FileRepository,
) *EvidenceService {
	return &EvidenceService{
		db:       db,
		repo:     repo,
		fileRepo: fileRepo,
	}
}

func (s *EvidenceService) Create(
	ctx context.Context,
	req dto.CreateEvidenceRequest,
) (*dto.EvidenceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	evidence := model.Evidence{
		ScheduleID: req.ScheduleID,
		FileID:     req.FileID,
	}

	res, err := s.repo.Create(ctx, tx, &evidence)
	if err != nil {
		return nil, err
	}

	file, err := s.fileRepo.FindByID(ctx, tx, evidence.FileID)
	if err != nil {
		return nil, err
	}

	res.File = dto.ToFileResponse(file)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *EvidenceService) GetByID(
	ctx context.Context,
	id uint,
) (*dto.EvidenceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	evidence, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if evidence == nil {
		return nil, customerrors.EvidenceErrNotFound
	}

	file, err := s.fileRepo.FindByID(ctx, tx, evidence.FileID)
	if err != nil {
		return nil, err
	}

	res := dto.ToEvidenceResponse(
		evidence,
		dto.ToFileResponse(file),
	)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *EvidenceService) GetByScheduleID(
	ctx context.Context,
	scheduleID uint,
) ([]*dto.EvidenceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	evidences, err := s.repo.FindByScheduleID(ctx, tx, scheduleID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.EvidenceResponse, 0, len(evidences))

	for _, evidence := range evidences {
		_, err := s.fileRepo.FindByID(ctx, tx, evidence.FileID)
		if err != nil {
			return nil, err
		}

		responses = append(
			responses,
			evidence,
		)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return responses, nil
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

	evidence, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if evidence == nil {
		return customerrors.EvidenceErrNotFound
	}

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}
