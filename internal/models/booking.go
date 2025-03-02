package models

import (
	"time"
)

type Booking struct {
	ID         string    `json:"id,omitempty"`
	ResourceID string    `json:"resource-id,omitempty"`
	UserID     int       `json:"user-id,omitempty"`
	StartTime  time.Time `json:"start-time"`
	EndTime    time.Time `json:"end-time"`
}

// type UpdateBookingInput struct {
// 	StartTime *time.Time `json:"start-time"`
// 	EndTime   *time.Time `json:"end-time"`
// }

// func (r *UpdateBookingInput) Validate() error {
// 	if r.StartTime == nil && r.EndTime == nil {
// 		return errors.New("at least one of start or end time must be provided")
// 	}
// 	if r.StartTime != nil && r.EndTime != nil {
// 		if r.StartTime.After(*r.EndTime) {
// 			return errors.New("start time cannot be latter than end time")
// 		}
// 	}

// 	return nil
// }
