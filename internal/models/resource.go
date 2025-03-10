package models

import (
	"errors"
)

type Resource struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Type   string `db:"type" json:"type"`
	Status string `json:"status" binding:"required,oneof=available booked"`
}

type CreateResourceInput struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status string `json:"status" binding:"required,oneof=available booked"`
}

type UpdateResourceInput struct {
	Name   *string `json:"name"`
	Type   *string `json:"type"`
	Status *string `json:"status"`
}

func (r *UpdateResourceInput) Validate() error {
	if r.Name == nil && r.Type == nil && r.Status == nil {
		return errors.New("empty update data")
	}

	return nil
}
