package models

import "time"

type Booking struct {
	ID         string    `json:"id,omitempty"`
	ResourceID string    `json:"resource-id,omitempty"`
	UserID     int       `json:"user-id,omitempty"`
	StartTime  time.Time `json:"start-time"`
	EndTime    time.Time `json:"end-time"`
}

type CreateBookingInput struct {
	ResourceID string
	UserID     int
	StartTime  time.Time
	EndTime    time.Time
}
