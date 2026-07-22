package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"magnolia-test-backend/internal/constants"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
	"magnolia-test-backend/internal/worker"
	"time"
)

type WorkingScheduleService struct {
	db           *sql.DB
	repo         *repository.WorkingScheduleRepository
	outletRepo   *repository.OutletRepository
	evidenceRepo *repository.EvidenceRepository
	worker       *worker.Worker
}

func NewWorkingScheduleService(db *sql.DB, repo *repository.WorkingScheduleRepository, outletRepo *repository.OutletRepository, evidenceRepo *repository.EvidenceRepository, worker *worker.Worker) *WorkingScheduleService {
	return &WorkingScheduleService{
		db:           db,
		repo:         repo,
		outletRepo:   outletRepo,
		evidenceRepo: evidenceRepo,
		worker:       worker,
	}
}

func (s *WorkingScheduleService) Create(
	ctx context.Context,
	req dto.CreateWorkingScheduleRequest,
) (*dto.WorkingScheduleResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	scheduleDate, err := time.Parse(time.RFC3339, req.ScheduleDate)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule_date: %w", err)
	}

	schedule := model.WorkingSchedule{
		OutletID:      req.OutletID,
		SalesID:       req.SalesID,
		Address:       req.Address,
		ScheduleDate:  scheduleDate,
		CurrentStage:  req.CurrentStage,
		ExpectedStage: req.ExpectedStage,
		Note:          req.Note,
		SyncStatus:    constants.SyncStatusQueued,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	res, err := s.repo.Upsert(ctx, tx, &schedule)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	s.triggerSync(res.ScheduleID)

	return res, nil
}

func (s *WorkingScheduleService) GetByID(
	ctx context.Context,
	id uint,
) (*dto.WorkingScheduleResponse, error) {
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

	outlet, err := s.outletRepo.FindByID(ctx, tx, schedule.OutletID)
	if err != nil {
		return nil, err
	}
	outletResponse := dto.ToOutletResponse(outlet, []*dto.WorkingScheduleResponse{})

	evidences, err := s.evidenceRepo.FindByScheduleID(ctx, tx, schedule.ScheduleID)
	if err != nil {
		return nil, err
	}

	res := dto.ToWorkingScheduleResponse(
		schedule,
		evidences,
		outletResponse,
	)

	return res, tx.Commit()
}

func (s *WorkingScheduleService) GetAll(
	ctx context.Context,
) ([]*dto.WorkingScheduleResponse, error) {
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
) (*dto.WorkingScheduleResponse, error) {
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

	scheduleDate, err := time.Parse(time.RFC3339, req.ScheduleDate)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule_date: %w", err)
	}

	schedule.SalesID = req.SalesID
	schedule.Address = req.Address
	schedule.ScheduleDate = scheduleDate
	schedule.ExpectedStage = req.ExpectedStage
	schedule.Note = req.Note
	schedule.SyncStatus = constants.SyncStatusQueued
	schedule.UpdatedAt = time.Now()

	res, err := s.repo.Upsert(ctx, tx, schedule)

	if err := s.evidenceRepo.DeleteByScheduleID(ctx, tx, schedule.ScheduleID); err != nil {
		return nil, err
	}

	for _, fileID := range req.FileIDs {
		_, err := s.evidenceRepo.Create(ctx, tx, &model.Evidence{
			ScheduleID: schedule.ScheduleID,
			FileID:     fileID,
		})

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	s.triggerSync(res.ScheduleID)

	return res, nil
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

func (s *WorkingScheduleService) ChangeStage(
	ctx context.Context,
	id uint,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Find schedule
	schedule, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if schedule == nil {
		return customerrors.WorkingScheduleErrNotFound
	}

	// 2. Check expected stage
	if schedule.ExpectedStage == nil || *schedule.ExpectedStage == "" {
		return fmt.Errorf("expected stage is required")
	}

	// 3. Check evidence
	evidences, err := s.evidenceRepo.FindByScheduleID(
		ctx,
		tx,
		schedule.ScheduleID,
	)

	if err != nil {
		return err
	}

	if len(evidences) == 0 {
		return fmt.Errorf("cannot change stage without evidence")
	}

	now := time.Now()

	// 4. Update outlet stage
	err = s.outletRepo.UpdateStage(
		ctx,
		tx,
		schedule.OutletID,
		*schedule.ExpectedStage,
		now,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkingScheduleService) triggerSync(scheduleID uint) {
	if scheduleID == 0 {
		return
	}

	go func(id uint) {
		if err := s.worker.SyncWorkingSchedule(context.Background(), id); err != nil {
			log.Printf("sync misa failed for schedule %d: %v", id, err)
		}
	}(scheduleID)
}
