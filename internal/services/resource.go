package service

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repository"
)

type ResourceService struct {
	repository repository.Resource
}

func NewResourceService(repository repository.Resource) *ResourceService {
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
