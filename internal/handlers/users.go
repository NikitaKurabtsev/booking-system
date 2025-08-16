package handlers

import (
	"net/http"

	"github.com/NikitaKurabtsev/booking-system/internal/handlers/dto"
	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @Tags auth
// @Description create user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body dto.SignUpUserDTO true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var dtoUser dto.SignUpUserDTO
	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	domainUser := dto.ConvertUserToDomain(dtoUser)

	if err := domainUser.Validate(); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.service.CreateUser(c.Request.Context(), domainUser)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// @Summary SignIn
// @Tags auth
// @Description login user
// @ID login
// @Accept  json
// @Produce  json
// @Param input body dto.SignInUserDTO true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var dtoUser dto.SignInUserDTO
	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	domainUser := dto.ConvertUserToDomain(dtoUser)

	token, err := h.service.User.GenerateToken(c.Request.Context(), domainUser.Username, domainUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
