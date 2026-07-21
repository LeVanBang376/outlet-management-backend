package evidence

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
	req CreateEvidenceRequest,
) error {

	evidence := Evidence{
		ScheduleID:  req.ScheduleID,
		ObjectKey:   req.ObjectKey,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Size:        req.Size,
		CreatedAt:   time.Now(),
	}

	return s.repo.Create(
		ctx,
		&evidence,
	)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uint,
) (*Evidence, error) {

	evidence, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	if evidence == nil {
		return nil, ErrNotFound
	}

	return evidence, nil
}

func (s *Service) GetByScheduleID(
	ctx context.Context,
	scheduleID uint,
) ([]Evidence, error) {

	return s.repo.FindByScheduleID(
		ctx,
		scheduleID,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	id uint,
) error {

	evidence, err := s.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	if evidence == nil {
		return ErrNotFound
	}

	return s.repo.Delete(
		ctx,
		id,
	)
}
