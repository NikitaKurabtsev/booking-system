package models

import "time"

type BookingResponse struct {
	ID           int       `json:"id"`
	ResourceName string    `json:"resource_name"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}
