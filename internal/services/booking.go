package service

import (
	"context"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repository"
)

type BookingService struct {
	repository repository.Booking
}

func NewBookingService(repository repository.Booking) *BookingService {
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

func (s *BookingService) CreateBooking(ctx context.Context, input domain.Booking) (int, error) {
	
	return 0, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, bookingID int) error {
	return nil
}
