package repository

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
}

type Resource interface {
	GetAll(ctx context.Context) ([]models.Resource, error)
}

type Booking interface {
	Create(ctx context.Context, input models.Booking) (models.Booking, error)
	GetAll(ctx context.Context, userID int) ([]models.Booking, error)
	Delete(ctx context.Context, bookingID int) error
}

type Repository struct {
	User
	Resource
	Booking
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Resource: NewResourceRepository(db),
		Booking:  NewBookingRepository(db),
	}
}
