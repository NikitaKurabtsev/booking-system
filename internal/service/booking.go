package services

import "github.com/NikitaKurabtsev/booking-system/internal/repository"

type BookingService struct {
	repository repository.Booking
}

func NewBookingService(repository repository.Booking) *BookingService {
	return &BookingService{repository: repository}
}
