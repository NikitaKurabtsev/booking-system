package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/repositories"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	signKey  = "z?_Z(75@q-H:-u/l5vtJ#E48yf=4jk"
	tokenTTL = 6 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type UserService struct {
	repository repositories.User
}

func NewUserService(repository repositories.User) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (domain.User, error) {
	user, err := s.repository.GetUser(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.Authenticate(ctx, username, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		user.ID,
	})

	return token.SignedString([]byte(signKey))
}

// TODO: where to call this method?
// TODO: check signKey
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

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	err := user.Validate()
	if err != nil {
		return 0, err
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hashedPassword

	userID, err := s.repository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
