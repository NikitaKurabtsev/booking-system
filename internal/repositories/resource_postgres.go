package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jmoiron/sqlx"
)

type ResourcesPostgres struct {
	db *sqlx.DB
}

func NewResourcePostgres(db *sqlx.DB) *ResourcesPostgres {
	return &ResourcesPostgres{db: db}
}

func (r *ResourcesPostgres) Create(input models.CreateResourceInput) (models.Resource, error) {
	const rawQuery = `
		INSERT INTO %s (name, type, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, type, status
	`
	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	var createdResource models.Resource
	err := r.db.QueryRowx(query, input.Name, input.Type, input.Status).StructScan(&createdResource)
	if err != nil {
		return models.Resource{}, fmt.Errorf("failed to create resource :%w", err)
	}

	return createdResource, nil
}

func (r *ResourcesPostgres) GetAll() ([]models.Resource, error) {
	const rawQuery = `
		SELECT id, name, type, status 
		FROM %s
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	var rescources []models.Resource
	err := r.db.Select(&rescources, query)
	if err != nil {
		return nil, err
	}

	return rescources, nil
}

func (r *ResourcesPostgres) GetById(resourceID int) (models.Resource, error) {
	const rawQuery = `
		SELECT id, name, type, status 
		FROM %s WHERE id = $1	
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	var resource models.Resource
	err := r.db.Get(&resource, query, resourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Resource{}, fmt.Errorf("resource with ID :%d does not exsit :%w", resourceID, err)
		}
		return models.Resource{}, fmt.Errorf("database error :%w", err)
	}

	return resource, nil
}

func (r *ResourcesPostgres) Update(resourceID int, input models.Resource) (models.Resource, error) {
	const rawQuery = `
		SELECT id, name, type, status
		FROM %s
		WHERE id = $1
	`
	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	var resource models.Resource
	err := r.db.Get(&resource, query, resourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Resource{}, fmt.Errorf("resource with ID: %d does not exist: %w", resourceID, err)
		}
		return models.Resource{}, err
	}

	//...
	const updateQuery = `
		UPDATE %s
		SET name = $1, type = $2, status = $3
		WHERE id = $4
		RETURNING id, name, type, status
	`

	query = fmt.Sprintf(updateQuery, db.ResourcesTable)

	var updatedResource models.Resource
	err = r.db.QueryRowx(query, input.Name, input.Type, input.Status).StructScan(&updatedResource)
	if err != nil {
		return models.Resource{}, err
	}

	return updatedResource, nil
}

func (r *ResourcesPostgres) Delete(resourceID int) error {
	const rawQuery = `
		DELETE FROM %s
		WHERE id = $1
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	_, err := r.db.Exec(query, resourceID)

	return err
}
