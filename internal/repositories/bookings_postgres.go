package repositories

import (
	"fmt"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jmoiron/sqlx"
)

type BookingPostgres struct {
	db *sqlx.DB
}

func NewBookingPostgres(db *sqlx.DB) *BookingPostgres {
	return &BookingPostgres{db: db}
}

func (r *BookingPostgres) Create(input models.Booking) (models.Booking, error) {
	rawQuery := `
		INSERT INTO %s (resource_id, user_id, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING id, resource_id, start_time, end_time 
	`

	query := fmt.Sprintf(rawQuery, db.BookingsTable)

	var booking models.Booking
	err := r.db.QueryRowx(query, input.ResourceID, input.UserID, input.StartTime, input.EndTime).StructScan(booking)
	if err != nil {
		return models.Booking{}, err
	}

	return booking, nil
}

func (r *BookingPostgres) GetAll(userID int) ([]models.Booking, error) {
	rawQuery := `
		SELECT b.id, r.name, b.start_time, b.end_time
		FROM %s b
		JOIN %s r ON b.resource_id = r.id
		WHERE b.user_id = $1
	`

	query := fmt.Sprintf(rawQuery, db.BookingsTable, db.ResourcesTable)

	var bookings []models.Booking
	err := r.db.Select(&bookings, query, userID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *BookingPostgres) Delete(bookingID int) error {
	rawQuery := `
		DELETE FROM %s
		WHERE id = $1
	`

	query := fmt.Sprintf(rawQuery, db.BookingsTable)

	_, err := r.db.Exec(query, bookingID)
	if err != nil {
		return err
	}

	return nil
}
