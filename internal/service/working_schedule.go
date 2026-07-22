package service

import (
	"context"
	"database/sql"
	"fmt"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
	"time"
)

type WorkingScheduleService struct {
	db   *sql.DB
	repo *repository.WorkingScheduleRepository
}

func NewWorkingScheduleService(db *sql.DB, repo *repository.WorkingScheduleRepository) *WorkingScheduleService {
	return &WorkingScheduleService{
		db:   db,
		repo: repo,
	}
}

func (s *WorkingScheduleService) Create(
	ctx context.Context,
	req dto.CreateWorkingScheduleRequest,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	scheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		return fmt.Errorf("invalid schedule_date: %w", err)
	}

	schedule := model.WorkingSchedule{
		OutletID:      req.OutletID,
		SalesID:       req.SalesID,
		Address:       req.Address,
		ScheduleDate:  scheduleDate,
		CurrentStage:  req.CurrentStage,
		ExpectedStage: req.ExpectedStage,
		Note:          req.Note,
		SyncStatus:    req.SyncStatus,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.repo.Create(ctx, tx, &schedule)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkingScheduleService) GetByID(
	ctx context.Context,
	id uint,
) (*model.WorkingSchedule, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	schedule, err := s.repo.FindByID(ctx, tx, id)

	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, customerrors.WorkingScheduleErrNotFound
	}

	return schedule, tx.Commit()
}

func (s *WorkingScheduleService) GetAll(
	ctx context.Context,
) ([]model.WorkingSchedule, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	schedules, err := s.repo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return schedules, tx.Commit()
}

func (s *WorkingScheduleService) Update(
	ctx context.Context,
	id uint,
	req dto.UpdateWorkingScheduleRequest,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	schedule, err := s.repo.FindByID(ctx, tx, id)

	if err != nil {
		return err
	}

	if schedule == nil {
		return customerrors.WorkingScheduleErrNotFound
	}

	scheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		return fmt.Errorf("invalid schedule_date: %w", err)
	}

	schedule.Address = req.Address
	schedule.ScheduleDate = scheduleDate
	schedule.CurrentStage = req.CurrentStage
	schedule.ExpectedStage = req.ExpectedStage
	schedule.Note = req.Note
	schedule.SyncStatus = req.SyncStatus
	schedule.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tx, schedule); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkingScheduleService) Delete(
	ctx context.Context,
	id uint,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}
