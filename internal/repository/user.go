package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, input models.User) (int, error) {
	rawQuery := `
		INSERT INTO %s (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	query := fmt.Sprintf(rawQuery, db.UsersTable)

	var userID int
	err := r.db.QueryRow(ctx, query, input.Username, input.Email, input.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	rawQuery := `
		SELECT id
		FROM %s
		WHERE username = $1 AND password_hash = $2
	`
	query := fmt.Sprintf(rawQuery, db.UsersTable)

	var user models.User
	err := r.db.QueryRow(ctx, query, username, password).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("user with Username :%s does not exsit :%w", username, err)
		}
		return models.User{}, err
	}

	return user, nil
}
