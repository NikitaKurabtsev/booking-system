package repositories

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/stretchr/testify/assert"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"

	"github.com/pashagolub/pgxmock/v2"
)

const (
	username = "test"
	email    = "test@gmail.com"
	password = "password"
	userID   = 1
)

func TestUserRepository_Create(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error occured when opening mock gpx pool: '%s'", err)
	}
	defer mockPool.Close()

	r := NewUserRepository(mockPool)

	tests := []struct {
		name    string
		mock    func()
		input   domain.User
		want    int
		wantErr bool
		errText error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				mockPool.ExpectBegin()
				mockPool.ExpectQuery("INSERT INTO users").
					WithArgs(username, email, password).
					WillReturnRows(rows)
				mockPool.ExpectCommit()
			},
			input: domain.User{
				Username: username,
				Email:    email,
				Password: password,
			},
			want: 1,
		},
		{
			name: "Unique Email Violation",
			mock: func() {
				mockPool.ExpectBegin()

				query := fmt.Sprintf("INSERT INTO %s", db.UsersTable)

				mockPool.ExpectQuery(query).
					WithArgs(username, email, password).
					WillReturnError(&pgconn.PgError{
						Code:           "23505",
						ConstraintName: "unique_email",
					})
				mockPool.ExpectRollback()
			},
			input: domain.User{
				Username: username,
				Email:    email,
				Password: password,
			},
			want:    0,
			wantErr: true,
			errText: ErrUserExists,
		},
		{
			name: "Unique Username Violation",
			mock: func() {
				mockPool.ExpectBegin()

				query := fmt.Sprintf("INSERT INTO %s", db.UsersTable)

				mockPool.ExpectQuery(query).
					WithArgs(username, email, password).
					WillReturnError(&pgconn.PgError{
						Code:           "23505",
						ConstraintName: "unique_username",
					})
				mockPool.ExpectRollback()
			},
			input: domain.User{
				Username: username,
				Email:    email,
				Password: password,
			},
			want:    0,
			wantErr: true,
			errText: ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, tt.errText, ErrUserExists)

			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	getQuery := `
		SELECT id, username, email, password
		FROM users
		WHERE username = $1
`
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error occured when opening mock gpx pool: '%s'", err)
	}
	defer mockPool.Close()

	r := NewUserRepository(mockPool)

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    domain.User
		wantErr bool
		errText string
	}{
		{
			name: "Ok",
			mock: func() {
				// TODO: hash???????????
				rows := pgxmock.NewRows([]string{"id", "username", "email", "password_hash"}).
					AddRow(1, "test", "test@gmail.com", "hashed_password")

				mockPool.ExpectQuery(regexp.QuoteMeta(getQuery)).
					WithArgs("test").
					WillReturnRows(rows)

			},
			input: "test",
			want: domain.User{
				ID: 1, Username: "test", Email: "test@gmail.com", Password: "hashed_password",
			},
		},
		{
			name: "User Not Found",
			mock: func() {
				mockPool.ExpectQuery(regexp.QuoteMeta(getQuery)).
					WithArgs("unknown").
					WillReturnError(pgx.ErrNoRows)
			},

			input:   "unknown",
			wantErr: true,
			errText: "user with Username :unknown does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)

				if len(tt.errText) != 0 {
					assert.Contains(t, err.Error(), tt.errText)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}

// func TestUserRepository_IsEmailExists(t *testing.T) {
// 	mockPool, err := pgxmock.NewPool()
// 	if err != nil {
// 		t.Fatalf("an error occured when opening mock gpx pool: '%s'", err)
// 	}
// 	defer mockPool.Close()

// 	r := NewUserRepository(mockPool)

// 	tests := []struct {
// 		name    string
// 		mock    func()
// 		input   string
// 		want    bool
// 		wantErr bool
// 		errText string
// 	}{
// 		{
// 			name: "Exists",
// 			mock: func() {
// 				rawQuery := `
// 					SELECT EXISTS(
// 					SELECT 1
// 					FROM %s
// 					WHERE email = $1
// 					)
// 					`

// 				query := fmt.Sprintf(rawQuery, db.UsersTable)

// 				rows := pgxmock.NewRows([]string{"exists"}).
// 					AddRow(true)

// 				mockPool.ExpectQuery(regexp.QuoteMeta(query)).
// 					WithArgs("test@gmail.com").
// 					WillReturnRows(rows)
// 			},
// 			input: "test@gmail.com",
// 			want:  true,
// 		},
// 		{
// 			name: "Not Exists",
// 			mock: func() {
// 				rawQuery := `
// 					SELECT EXISTS(
// 					SELECT 1
// 					FROM %s
// 					WHERE email = $1
// 					)
// 					`

// 				query := fmt.Sprintf(rawQuery, db.UsersTable)

// 				rows := pgxmock.NewRows([]string{"exists"}).
// 					AddRow(false)

// 				mockPool.ExpectQuery(regexp.QuoteMeta(query)).
// 					WithArgs("test@gmail.com").
// 					WillReturnRows(rows)
// 			},
// 			input: "test@gmail.com",
// 			want:  false,
// 		},
// 		{
// 			name: "Database Error",
// 			mock: func() {
// 				rawQuery := `
// 					SELECT EXISTS(
// 					SELECT 1
// 					FROM %s
// 					WHERE email = $1
// 					)
// 					`

// 				query := fmt.Sprintf(rawQuery, db.UsersTable)

// 				mockPool.ExpectQuery(regexp.QuoteMeta(query)).
// 					WithArgs("test@gmail.com").
// 					WillReturnError(errors.New("failed to check email exsitstance"))
// 			},
// 			input:   "test@gmail.com",
// 			want:    false,
// 			wantErr: true,
// 			errText: "failed to check email exsitstance",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()

// 			got, err := r.IsEmailExists(context.Background(), tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)

// 				if len(tt.errText) != 0 {
// 					assert.Contains(t, err.Error(), tt.errText)
// 				} else {
// 					assert.NoError(t, err)
// 					assert.Equal(t, tt.want, got)
// 				}
// 			}

// 			assert.NoError(t, mockPool.ExpectationsWereMet())
// 		})
// 	}
// }
