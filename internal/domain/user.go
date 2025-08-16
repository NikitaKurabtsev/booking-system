package domain

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email")

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

func (u *User) Validate() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}
