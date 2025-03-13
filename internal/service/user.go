package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/NikitaKurabtsev/booking-system/internal/models"
	"github.com/NikitaKurabtsev/booking-system/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"time"
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
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GenerateToken(username, password string) (string, error) {
	user, err := s.repository.GetUser(username, hashPassword(password))
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

func (s *UserService) CreateUser(user models.User) (int, error) {
	user.PasswordHash = hashPassword(user.PasswordHash)

	return s.repository.Create(user)
}

func hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
