package models

import (
	"errors"
	"time"
)

type Booking struct {
	ID         string    `json:"id,omitempty" db:"id"`
	UserID     int       `json:"user_id,omitempty" db:"user_id"`
	ResourceID string    `json:"resource_id,omitempty" db:"resource_id"`
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
}

type UpdateBookingInput struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

func (r *UpdateBookingInput) Validate() error {
	if r.StartTime == nil && r.EndTime == nil {
		return errors.New("at least one of start or end time must be provided")
	}
	if r.StartTime != nil && r.EndTime != nil {
		if r.StartTime.After(*r.EndTime) {
			return errors.New("start time cannot be latter than end time")
		}
	}

	return nil
}
