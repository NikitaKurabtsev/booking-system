package handlers

import (
	"net/http"

	"github.com/NikitaKurabtsev/booking-system/internal/handlers/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getBookings(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := h.service.Booking.GetBookings(c.Request.Context(), userID)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// @Summary Create booking
// @Security ApiKeyAuth
// @Tags bookings
// @Description create booking
// @ID create-booking
// @Accept  json
// @Produce  json
// @Param input body dto.CreateBookingDTO true "booking info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/bookings [post]
func (h *Handler) createBooking(c *gin.Context) {
	var booking dto.CreateBookingDTO

	userID, err := getUserID(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&booking); err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	booking.UserID = userID

	domainBooking, err := dto.ConvertBookingToDomain(booking)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bookingID, err := h.service.Booking.CreateBooking(c.Request.Context(), domainBooking)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return

	}

	c.JSON(http.StatusCreated, gin.H{"booking_id": bookingID})
}

func (h *Handler) deleteBooking(c *gin.Context) {
	bookingID, err := getBookingID(c)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Booking.DeleteBooking(c, bookingID)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "booking delete successfully"})
}
