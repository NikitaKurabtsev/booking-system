package dto

import (
	"errors"
	"time"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
)

var ErrInvalidBookingDateTime = errors.New("date should be like 2006-01-02 15:04:05")

type CreateBookingDTO struct {
	UserID     int    `json:"user_id,omitempty"`
	ResourceID int    `json:"resource_id" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

func ConvertBookingToDomain(dto CreateBookingDTO) (domain.Booking, error) {
	const layout = "2006-01-02 15:04:05"

	startTime, err := time.Parse(layout, dto.StartTime)
	if err != nil {
		return domain.Booking{}, ErrInvalidBookingDateTime
	}

	endTime, err := time.Parse(layout, dto.EndTime)
	if err != nil {
		return domain.Booking{}, ErrInvalidBookingDateTime
	}

	return domain.Booking{
		ResourceID: dto.ResourceID,
		UserID:     dto.UserID,
		StartTime:  startTime,
		EndTime:    endTime,
	}, nil
}
