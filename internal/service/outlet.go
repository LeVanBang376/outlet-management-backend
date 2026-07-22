package service

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/constants"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
	"magnolia-test-backend/internal/response"
	"time"
)

type OutletService struct {
	db                  *sql.DB
	repo                *repository.OutletRepository
	workingScheduleRepo *repository.WorkingScheduleRepository
}

func NewOutletService(db *sql.DB, repo *repository.OutletRepository, wsRepo *repository.WorkingScheduleRepository) *OutletService {
	return &OutletService{
		db:                  db,
		repo:                repo,
		workingScheduleRepo: wsRepo,
	}
}

func (s *OutletService) Create(ctx context.Context, req dto.CreateOutletRequest) (*dto.OutletResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	outlet := model.Outlet{
		Name:      req.Name,
		Address:   req.Address,
		Channel:   req.Channel,
		Tier:      req.Tier,
		SalesID:   req.SalesID,
		Stage:     req.Stage,
		Note:      &req.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res, err := s.repo.Create(ctx, tx, &outlet)
	if err != nil {
		return nil, err
	}

	if req.HasWorkingSchedule {
		schedule := &model.WorkingSchedule{
			OutletID:      res.OutletID,
			SalesID:       res.SalesID,
			Address:       res.Address,
			ScheduleDate:  now,
			CurrentStage:  res.Stage,
			ExpectedStage: nil,
			Note:          nil,
			SyncStatus:    constants.SyncStatusSynced,
		}

		err := s.workingScheduleRepo.Create(ctx, tx, schedule)
		if err != nil {
			return nil, err
		}
	}

	return res, tx.Commit()
}

func (s *OutletService) GetByID(ctx context.Context, id uint) (*model.Outlet, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	outlet, err := s.repo.FindByID(ctx, tx, id)

	if err != nil {
		return nil, err
	}

	if outlet == nil {
		return nil, customerrors.OutletErrNotFound
	}

	return outlet, tx.Commit()
}

func (s *OutletService) GetAll(ctx context.Context, pagination *response.Pagination) ([]model.Outlet, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	outlets, err := s.repo.FindAll(ctx, tx, pagination)
	if err != nil {
		return nil, err
	}

	return outlets, tx.Commit()
}

func (s *OutletService) Update(ctx context.Context, id uint, req dto.UpdateOutletRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	outlet, err := s.repo.FindByID(ctx, tx, id)

	if err != nil {
		return err
	}

	if outlet == nil {
		return customerrors.OutletErrNotFound
	}

	outlet.Name = req.Name
	outlet.Address = req.Address
	outlet.Channel = req.Channel
	outlet.Tier = req.Tier
	outlet.SalesID = req.SalesID
	outlet.Stage = req.Stage
	outlet.Note = &req.Note
	outlet.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tx, outlet); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OutletService) Delete(ctx context.Context, id uint) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	outlet, err := s.repo.FindByID(ctx, tx, id)

	if err != nil {
		return err
	}

	if outlet == nil {
		return customerrors.OutletErrNotFound
	}

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}
