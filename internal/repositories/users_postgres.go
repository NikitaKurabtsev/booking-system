package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(input models.User) (int, error) {
	rawQuery := `
		INSERT INTO %s (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	query := fmt.Sprintf(rawQuery, db.UsersTable)

	var userID int
	err := r.db.QueryRowx(query, input.Username, input.Email, input.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserPostgres) GetUser(username, password string) (models.User, error) {
	rawQuery := `
		SELECT id
		FROM %s
		WHERE username = $1 AND password_hash = $2
	`
	query := fmt.Sprintf(rawQuery, db.UsersTable)

	var user models.User
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("user with Username :%s does not exsit :%w", username, err)
		}
		return models.User{}, err
	}

	return user, nil
}
