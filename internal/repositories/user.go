package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type user struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserRepository struct {
	db DBPool
}

func NewUserRepository(db DBPool) *UserRepository {
	return &UserRepository{db: db}
}

func toDomainUser(u user) domain.User {
	return domain.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func toDBUser(u domain.User) user {
	return user{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (r *UserRepository) Create(ctx context.Context, inputUser domain.User) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil && !errors.Is(txErr, pgx.ErrTxClosed) {
				log.Printf("failed to rollback transaction %v", txErr)

			}
		}
	}()

	dbUser := toDBUser(inputUser)

	query := `
		INSERT INTO users 
			(username, email, password)
		VALUES 
			($1, $2, $3)
		RETURNING id
	`

	var userID int
	err = tx.QueryRow(ctx, query, dbUser.Username, dbUser.Email, dbUser.Password).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("user already exsits: %v", err)
			if pgErr.ConstraintName == "unique_username" {
				return 0, fmt.Errorf("provided username is already exists: %w", ErrUserExists)
			}
			if pgErr.ConstraintName == "unique_email" {
				return 0, fmt.Errorf("provided email is already exists: %w", ErrUserExists)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username string) (domain.User, error) {
	query := `
		SELECT id, username, email, password
		FROM users
		WHERE username = $1;
	`

	var dbUser user
	err := r.db.QueryRow(ctx, query, username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with Username :%s does not exist :%w", username, err)
		}
		return domain.User{}, err
	}

	domainUser := toDomainUser(dbUser)

	return domainUser, nil
}
