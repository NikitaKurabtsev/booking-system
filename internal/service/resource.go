package services

import "github.com/NikitaKurabtsev/booking-system/internal/repository"

type ResourceService struct {
	repository repository.Resource
}

func NewResourceService(repository repository.Resource) *ResourceService {
	return &ResourceService{
		repository: repository,
	}
}
