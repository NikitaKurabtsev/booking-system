package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jmoiron/sqlx"
)

type ResourcePostgres struct {
	db *sqlx.DB
}

func NewResourcePostgres(db *sqlx.DB) *ResourcePostgres {
	return &ResourcePostgres{db: db}
}

func (r *ResourcePostgres) Create(input models.CreateResourceInput) (models.Resource, error) {
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

func (r *ResourcePostgres) GetAll() ([]models.Resource, error) {
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

func (r *ResourcePostgres) GetById(resourceID int) (models.Resource, error) {
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

func (r *ResourcePostgres) Update(resourceID int, input models.UpdateResourceInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Type != nil {
		setValues = append(setValues, fmt.Sprintf("type=$%d", argId))
		args = append(args, *input.Type)
		argId++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	setArgs := strings.Join(setValues, ", ")

	const rawQuery = `
		UPDATE %s
		SET %s 
		WHERE id = %d
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable, setArgs, argId)

	args = append(args, resourceID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("resource with ID: %d does not exist: %w", resourceID, err)
		}
	}

	return nil
}

func (r *ResourcePostgres) Delete(resourceID int) error {
	const rawQuery = `
		DELETE FROM %s
		WHERE id = $1
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	_, err := r.db.Exec(query, resourceID)

	return err
}
