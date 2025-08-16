package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidBookingTime     = errors.New("end time must be after start time")
	ErrInvalidBookingDuration = errors.New("max booking duration is 4 hours")
)

type Booking struct {
	ID         int
	ResourceID int
	UserID     int
	StartTime  time.Time
	EndTime    time.Time
}
