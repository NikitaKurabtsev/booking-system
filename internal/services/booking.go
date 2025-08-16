package services

import (
	"context"
	"errors"
	"time"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repositories"
)

var (
	ErrInvalidBookingTime     = errors.New("end time must be after start time")
	ErrInvalidBookingDuration = errors.New("max booking duration is 4 hours")
)

type BookingService struct {
	repository repositories.Booking
}

func NewBookingService(repository repositories.Booking) *BookingService {
	return &BookingService{
		repository: repository,
	}
}

func (s *BookingService) GetBookings(ctx context.Context, userID int) ([]domain.Booking, error) {
	bookings, err := s.repository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *BookingService) CreateBooking(ctx context.Context, booking domain.Booking) (int, error) {
	if booking.StartTime.After(booking.EndTime) {
		return 0, ErrInvalidBookingTime
	}

	if booking.EndTime.Before(booking.StartTime) {
		return 0, ErrInvalidBookingTime
	}

	if booking.EndTime.Sub(booking.StartTime) > 4*time.Hour {
		return 0, ErrInvalidBookingDuration
	}

	bookingID, err := s.repository.Create(ctx, booking)
	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, bookingID int) error {
	err := s.repository.Delete(ctx, bookingID)
	if err != nil {
		return err
	}

	return nil
}
