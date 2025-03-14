package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

const (
	salt     = "K3o)u?wb*SLvBQ-n39GH"
	signKey  = "z?_Z(75@q-H:-u/l5vtJ#E48yf=4jk"
	tokenTTL = 6 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type UserService struct {
	repository repository.User
	logger     *slog.Logger
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	hashedPassword := hashPassword(password)

	user, err := s.repository.GetUser(ctx, username, hashedPassword)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.ID,
	})

	return token.SignedString([]byte(signKey))
}

func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("expected *tokenClaims, got different type")
	}

	return claims.UserID, nil
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (int, error) {
	err := isValidEmail(user.Email)
	if err != nil {
		return 0, err
	}

	hashedPassword := hashPassword(user.Password)
	user.Password = hashedPassword

	return s.repository.Create(ctx, user)
}

func isValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email")
	}

	return nil
}

func hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
