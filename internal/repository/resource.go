package repository

import (
	"context"
	"fmt"
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceRepository struct {
	db *pgxpool.Pool
}

func NewResourceRepository(db *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{db: db}
}

//func (r *ResourceRepository) Create(input models.CreateResourceInput) (models.Resource, error) {
//	const rawQuery = `
//		INSERT INTO %s (name, type, status)
//		VALUES ($1, $2, $3)
//		RETURNING id, name, type, status
//	`
//	query := fmt.Sprintf(rawQuery, db.ResourcesTable)
//
//	var createdResource models.Resource
//	err := r.db.QueryRowx(query, input.Name, input.Type, input.Status).StructScan(&createdResource)
//	if err != nil {
//		return models.Resource{}, fmt.Errorf("failed to create resource :%w", err)
//	}
//
//	return createdResource, nil
//}

func (r *ResourceRepository) GetAll(ctx context.Context) ([]models.Resource, error) {
	const rawQuery = `
		SELECT id, name, type, status 
		FROM %s
	`

	query := fmt.Sprintf(rawQuery, db.ResourcesTable)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []models.Resource

	for rows.Next() {
		var resource models.Resource
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Status, &resource.Type)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil

	//err := r.db.Select(&resources, query)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return resources, nil
}

//func (r *ResourceRepository) GetById(resourceID int) (models.Resource, error) {
//	const rawQuery = `
//		SELECT id, name, type, status
//		FROM %s WHERE id = $1
//	`
//
//	query := fmt.Sprintf(rawQuery, db.ResourcesTable)
//
//	var resource models.Resource
//	err := r.db.Get(&resource, query, resourceID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return models.Resource{}, fmt.Errorf("resource with ID :%d does not exsit :%w", resourceID, err)
//		}
//		return models.Resource{}, fmt.Errorf("database error :%w", err)
//	}
//
//	return resource, nil
//}

//func (r *ResourceRepository) Update(resourceID int, input models.UpdateResourceInput) error {
//	setValues := make([]string, 0)
//	args := make([]interface{}, 0)
//	argId := 1
//
//	if input.Name != nil {
//		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
//		args = append(args, *input.Name)
//		argId++
//	}
//
//	if input.Type != nil {
//		setValues = append(setValues, fmt.Sprintf("type=$%d", argId))
//		args = append(args, *input.Type)
//		argId++
//	}
//
//	if input.Status != nil {
//		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
//		args = append(args, *input.Status)
//		argId++
//	}
//
//	setArgs := strings.Join(setValues, ", ")
//
//	const rawQuery = `
//		UPDATE %s
//		SET %s
//		WHERE id = %d
//	`
//
//	query := fmt.Sprintf(rawQuery, db.ResourcesTable, setArgs, argId)
//
//	args = append(args, resourceID)
//
//	_, err := r.db.Exec(query, args...)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return fmt.Errorf("resource with ID: %d does not exist: %w", resourceID, err)
//		}
//	}
//
//	return nil
//}

//func (r *ResourceRepository) Delete(resourceID int) error {
//	const rawQuery = `
//		DELETE FROM %s
//		WHERE id = $1
//	`
//
//	query := fmt.Sprintf(rawQuery, db.ResourcesTable)
//
//	_, err := r.db.Exec(query, resourceID)
//
//	return err
//}
