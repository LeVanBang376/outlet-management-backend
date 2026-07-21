package outlet

import (
	"context"
	"magnolia-test-backend/internal/response"
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

func (s *Service) Create(ctx context.Context, req CreateOutletRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	now := time.Now()

	outlet := Outlet{
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

	return s.repo.Create(ctx, &outlet)
}

func (s *Service) GetByID(ctx context.Context, id uint) (*Outlet, error) {
	outlet, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if outlet == nil {
		return nil, ErrNotFound
	}

	return outlet, nil
}

func (s *Service) GetAll(ctx context.Context, pagination *response.Pagination) ([]Outlet, error) {
	return s.repo.FindAll(ctx, pagination)
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateOutletRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	outlet, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	if outlet == nil {
		return ErrNotFound
	}

	outlet.Name = req.Name
	outlet.Address = req.Address
	outlet.Channel = req.Channel
	outlet.Tier = req.Tier
	outlet.SalesID = req.SalesID
	outlet.Stage = req.Stage
	outlet.Note = &req.Note
	outlet.UpdatedAt = time.Now()

	return s.repo.Update(ctx, outlet)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	outlet, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	if outlet == nil {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
