package working_schedule

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateWorkingScheduleRequest,
) error {

	now := time.Now()

	schedule := WorkingSchedule{
		OutletID:      req.OutletID,
		SalesID:       req.SalesID,
		Address:       req.Address,
		ScheduleDate:  req.ScheduleDate,
		CurrentStage:  req.CurrentStage,
		ExpectedStage: req.ExpectedStage,
		Note:          req.Note,
		SyncStatus:    req.SyncStatus,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return s.repo.Create(ctx, &schedule)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uint,
) (*WorkingSchedule, error) {
	schedule, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, ErrNotFound
	}

	return schedule, nil
}

func (s *Service) GetAll(
	ctx context.Context,
) ([]WorkingSchedule, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) Update(
	ctx context.Context,
	id uint,
	req UpdateWorkingScheduleRequest,
) error {

	schedule, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	if schedule == nil {
		return ErrNotFound
	}

	schedule.Address = req.Address
	schedule.ScheduleDate = req.ScheduleDate
	schedule.CurrentStage = req.CurrentStage
	schedule.ExpectedStage = req.ExpectedStage
	schedule.Note = req.Note
	schedule.SyncStatus = req.SyncStatus
	schedule.UpdatedAt = time.Now()

	return s.repo.Update(ctx, schedule)
}

func (s *Service) Delete(
	ctx context.Context,
	id uint,
) error {
	_, err := s.GetByID(ctx, id)

	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
