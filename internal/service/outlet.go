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
	"magnolia-test-backend/internal/response"
	"magnolia-test-backend/internal/worker"
	"time"
)

type OutletService struct {
	db                  *sql.DB
	repo                *repository.OutletRepository
	workingScheduleRepo *repository.WorkingScheduleRepository
	worker              *worker.Worker
}

func NewOutletService(db *sql.DB, repo *repository.OutletRepository, wsRepo *repository.WorkingScheduleRepository, worker *worker.Worker) *OutletService {
	return &OutletService{
		db:                  db,
		repo:                repo,
		workingScheduleRepo: wsRepo,
		worker:              worker,
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

	var scheduleID uint

	if req.HasWorkingSchedule {
		scheduleDate, err := time.Parse(time.RFC3339, req.ScheduleDate)
		if err != nil {
			return nil, fmt.Errorf("invalid schedule_date: %w", err)
		}

		schedule := &model.WorkingSchedule{
			OutletID:      res.OutletID,
			SalesID:       res.SalesID,
			Address:       res.Address,
			ScheduleDate:  scheduleDate,
			CurrentStage:  res.Stage,
			ExpectedStage: &req.ExpectedStage,
			Note:          &req.ScheduleNote,
			SyncStatus:    constants.SyncStatusQueued,
		}

		scheduleRes, err := s.workingScheduleRepo.Upsert(ctx, tx, schedule)
		if err != nil {
			return nil, err
		}

		scheduleID = scheduleRes.ScheduleID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if scheduleID > 0 {
		go func(id uint) {
			if err := s.worker.SyncWorkingSchedule(context.Background(), id); err != nil {
				log.Printf("sync misa failed: %v", err)
			}
		}(scheduleID)
	}

	return res, nil
}

func (s *OutletService) GetByID(ctx context.Context, id uint) (*dto.OutletResponse, error) {
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

	schedules, err := s.workingScheduleRepo.FindByOutletID(ctx, tx, outlet.OutletID)
	if err != nil {
		return nil, err
	}

	schedulesResponse := make([]*dto.WorkingScheduleResponse, 0)
	for _, schedule := range schedules {
		schedulesResponse = append(schedulesResponse, dto.ToWorkingScheduleResponse(&schedule, []*dto.EvidenceResponse{}, nil))
	}

	res := dto.ToOutletResponse(outlet, schedulesResponse)

	return res, tx.Commit()
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
