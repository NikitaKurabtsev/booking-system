package services

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repositories"
)

type User interface {
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	CreateUser(ctx context.Context, user domain.User) (int, error)
	Authenticate(ctx context.Context, username, password string) (domain.User, error)
}

type Resource interface {
	GetResources(ctx context.Context) ([]domain.Resource, error)
}

type Booking interface {
	GetBookings(ctx context.Context, userID int) ([]domain.Booking, error)
	CreateBooking(ctx context.Context, booking domain.Booking) (int, error)
	DeleteBooking(ctx context.Context, bookingID int) error
}

type Service struct {
	User
	Resource
	Booking
}

func NewService(repositories *repositories.Repository) *Service {
	return &Service{
		User:     NewUserService(repositories.User),
		Resource: NewResourceService(repositories.Resource),
		Booking:  NewBookingService(repositories.Booking),
	}
}
