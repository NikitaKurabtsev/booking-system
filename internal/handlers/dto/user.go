package dto

import "github.com/NikitaKurabtsev/booking-system/internal/domain"

type userDTO interface{}

type SignUpUserDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// type UserResponse struct {
// ...
//}

func ConvertUserToDomain(dto userDTO) domain.User {
	switch v := dto.(type) {
	case SignUpUserDTO:
		return domain.User{
			Username: v.Username,
			Email:    v.Email,
			Password: v.Password,
		}
	case SignInUserDTO:
		return domain.User{
			Username: v.Username,
			Email:    "",
			Password: v.Password,
		}
	default:
		return domain.User{}
	}
}
