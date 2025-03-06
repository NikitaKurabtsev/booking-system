package repositories

import (
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type ResourceRepository interface {
	// Create(inputResource models.Resource) (models.Resource, error)
	GetAll() ([]models.Resource, error)
	// GetById(resourceID int) (models.Resource, error)
	// Update(inputResource models.Resource) (models.Resource, error)
	// Delete(resourceID int) error
}

type BookingRepository interface {
	Create(input models.Booking) (models.Booking, error)
	GetAll(userID int) ([]models.Booking, error)
	Delete(bookingID int) error
}

type Repository struct {
	UserRepository
	ResourceRepository
	BookingRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:     NewUserPostgres(db),
		ResourceRepository: NewResourcePostgres(db),
		BookingRepository:  NewBookingPostgres(db),
	}
}
