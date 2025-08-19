package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

type booking struct {
	ID         int
	UserID     int
	ResourceID int
	StartTime  pgtype.Timestamp
	EndTime    pgtype.Timestamp
}

type BookingRepository struct {
	db    DBPool
	cache cache.Cache
}

func NewBookingRepository(db DBPool, cache cache.Cache) *BookingRepository {
	return &BookingRepository{
		db:    db,
		cache: cache,
	}
}

func toDomainBooking(b booking) domain.Booking {
	return domain.Booking{
		ID:         b.ID,
		UserID:     b.UserID,
		ResourceID: b.ResourceID,
		StartTime:  b.StartTime.Time,
		EndTime:    b.EndTime.Time,
	}
}

func toDBBooking(b domain.Booking) booking {
	return booking{
		ID:         b.ID,
		UserID:     b.UserID,
		ResourceID: b.ResourceID,
		StartTime:  pgtype.Timestamp{Time: b.StartTime.UTC(), Valid: !b.StartTime.IsZero()},
		EndTime:    pgtype.Timestamp{Time: b.EndTime.UTC(), Valid: !b.EndTime.IsZero()},
	}
}

func (r *BookingRepository) Create(ctx context.Context, inputBooking domain.Booking) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil && !errors.Is(txErr, pgx.ErrTxClosed) {
				log.Printf("failed to rollback transaction: %v", txErr)

			}
		}
	}()

	dbBooking := toDBBooking(inputBooking)

	//check overlap

	// TODO: check query and args
	query := `
		SELECT start_time, end_time
		FROM bookings
		WHERE resource_id = $1
		AND $2 < end_time
		AND $3 > start_time
`

	var overlapStart, overlapEnd time.Time
	err = tx.QueryRow(
		ctx,
		query,
		dbBooking.ResourceID,
		dbBooking.StartTime.Time,
		dbBooking.EndTime.Time,
	).Scan(&overlapStart, &overlapEnd)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("failed to check overlapping bookings: %w", err)
	}
	if err == nil {
		return 0, fmt.Errorf(
			"resource is already booked from %s to %s",
			overlapStart.Format(time.RFC3339),
			overlapEnd.Format(time.RFC3339),
		)
	}

	query = `
		INSERT INTO bookings
			(resource_id, user_id, start_time, end_time)
		VALUES
			($1, $2, $3, $4)
		RETURNING id
`

	var bookingID int
	err = tx.QueryRow(
		ctx,
		query,
		dbBooking.ResourceID,
		dbBooking.UserID,
		dbBooking.StartTime,
		dbBooking.EndTime,
	).Scan(&bookingID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		if err := r.cache.Delete(ctx, "resources:all"); err != nil {
			log.Printf("failed to invalidate cache: %v", err)
		}
	}()

	return bookingID, nil
}

// check where to get user id from params or from the headers
func (r *BookingRepository) GetAll(ctx context.Context, userID int) ([]domain.Booking, error) {
	// TODO: where to get cache?????
	cacheKey := fmt.Sprintf("bookings:users:%d", userID) // TODO: check cache key

	cached, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var bookings []domain.Booking
		if err = json.Unmarshal([]byte(cached), &bookings); err != nil {
			return nil, fmt.Errorf("cache unmarshal error: %w", err)
		}
		return bookings, nil
	} else {
		log.Printf("cache GET error: %v", err)
	}

	query := `
		SELECT id, resource_id, start_time, end_time
		FROM bookings
		WHERE user_id = $1
`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to scan a row: %w", err)
	}
	defer rows.Close()

	var bookings []booking
	for rows.Next() {
		var b booking
		err = rows.Scan(&b.ID, &b.ResourceID, &b.StartTime, &b.EndTime)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	domainBookings := make([]domain.Booking, 0, len(bookings))
	for _, b := range bookings {
		domainBookings = append(domainBookings, toDomainBooking(b))
	}

	go func() {
		jsonData, err := json.Marshal(domainBookings)
		if err != nil {
			log.Printf("marshal error: %v", err)
		}

		err = r.cache.Set(context.Background(), cacheKey, jsonData, cache.TTL)
		if err != nil {
			log.Printf("cache SET error: %v", err)
		}
	}()

	return domainBookings, nil
}

func (r *BookingRepository) Delete(ctx context.Context, bookingID int) error {
	query := `
		DELETE bookings
		WHERE id = $1
`

	affectedRows, err := r.db.Exec(ctx, query, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if affectedRows.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", bookingID)
	}

	return nil
}
