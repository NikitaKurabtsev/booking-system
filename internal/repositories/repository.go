package repositories

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/pkg/cache"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type User interface {
	Create(ctx context.Context, inputUser domain.User) (int, error)
	GetUser(ctx context.Context, username string) (domain.User, error)
}

type Resource interface {
	GetAll(ctx context.Context) ([]domain.Resource, error)
}

type Booking interface {
	Create(ctx context.Context, inputBooking domain.Booking) (int, error)
	GetAll(ctx context.Context, userID int) ([]domain.Booking, error)
	Delete(ctx context.Context, bookingID int) error
}

type DBPool interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

// type Tx interface {
// 	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
// 	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
// 	Commit(ctx context.Context) error
// 	Rollback(ctx context.Context) error
// }

type Repository struct {
	User
	Resource
	Booking
}

func NewRepository(db DBPool, cache cache.Cache) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Resource: NewResourceRepository(db, cache),
		Booking:  NewBookingRepository(db, cache),
	}
}
