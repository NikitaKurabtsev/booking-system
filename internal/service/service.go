package services

import (
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/internal/repository"
)

type User interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	CreateUser(user models.User) (int, error)
}

type Resource interface {
	// TODO: add methods
	//GetAll(userID int) ([]models.Resource, error)
}

type Booking interface {
	// TODO: add methods
}

type Service struct {
	User
	Resource
	Booking
}

func NewService(repositories repository.Repository) *Service {
	return &Service{
		User:     NewUserService(repositories.User),
		Resource: NewResourceService(repositories.Resource),
		Booking:  NewBookingService(repositories.Booking),
	}
}
