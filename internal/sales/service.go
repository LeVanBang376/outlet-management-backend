package sales

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]SalesResponse, error) {
	sales, err := s.repository.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]SalesResponse, 0, len(sales))

	for _, item := range sales {
		result = append(result, NewSalesResponse(item))
	}

	return result, nil
}
