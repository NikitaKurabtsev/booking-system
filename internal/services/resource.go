package services

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repositories"
)

type ResourceService struct {
	repository repositories.Resource
}

func NewResourceService(repository repositories.Resource) *ResourceService {
	return &ResourceService{
		repository: repository,
	}
}

func (s *ResourceService) GetResources(ctx context.Context) ([]domain.Resource, error) {
	resources, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
