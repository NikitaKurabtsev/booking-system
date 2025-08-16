package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/NikitaKurabtsev/booking-system/internal/handlers/dto"
	"github.com/NikitaKurabtsev/booking-system/internal/services"
	mocksbookingservice "github.com/NikitaKurabtsev/booking-system/internal/services/mocks"
	mocklogger "github.com/NikitaKurabtsev/booking-system/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_createBooking(t *testing.T) {
	type mockBehavior func(s *mocksbookingservice.MockBooking, logger *mocklogger.MockLogger)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		inputData            dto.CreateBookingDTO
		requestMethod        string
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mocksbookingservice.MockBooking, l *mocklogger.MockLogger) {
				s.EXPECT().
					CreateBooking(gomock.Any(), gomock.Any()).
					Return(1, nil)
			},
			inputData: dto.CreateBookingDTO{
				ResourceID: 1,
				UserID:     0,
				StartTime:  "2025-11-01 13:00:00",
				EndTime:    "2025-11-01 14:00:00",
			},
			requestMethod:        "POST",
			expectedCodeStatus:   201,
			expectedResponseBody: `{"booking_id":1"}`,
		},
		{
			name: "Bad Request",
			mockBehavior: func(s *mocksbookingservice.MockBooking, l *mocklogger.MockLogger) {
				s.EXPECT().
					CreateBooking(gomock.Any(), gomock.Any()).
					Return(0, errors.New("something went wrong"))
				l.EXPECT().Error("error", "something went wrong")
			},
			inputData: dto.CreateBookingDTO{
				ResourceID: 1,
				UserID:     0,
				StartTime:  "2023-11-01 13:00:00",
				EndTime:    "2025-11-01 14:00:00",
			},
			requestMethod:        "POST",
			expectedCodeStatus:   400,
			expectedResponseBody: `{"error":""something went wrong"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLogger := mocklogger.NewMockLogger(ctrl)
			mockBookingService := mocksbookingservice.NewMockBooking(ctrl)

			test.mockBehavior(mockBookingService, mockLogger)

			bookingService := &services.Service{Booking: mockBookingService}
			handler := &Handler{service: bookingService, logger: mockLogger}

			w := httptest.NewRecorder()

			jsonData, _ := json.Marshal(test.inputData)
			req := httptest.NewRequest(test.requestMethod, "/booking", bytes.NewReader(jsonData))

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userID", 1)
				c.Next()
			})
			router.Handle(test.requestMethod, "/booking", handler.createBooking)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCodeStatus, w.Code)

		})
	}
}

func TestHandler_deleteBooking(t *testing.T) {
	type mockBehavior func(s *mocksbookingservice.MockBooking, logger *mocklogger.MockLogger)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		requestMethod        string
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mocksbookingservice.MockBooking, l *mocklogger.MockLogger) {
				s.EXPECT().
					DeleteBooking(gomock.Any(), 1).
					Return(nil)

			},
			requestMethod:        "DELETE",
			expectedCodeStatus:   204,
			expectedResponseBody: ``,
		},
		{
			name: "Bad Request",
			mockBehavior: func(s *mocksbookingservice.MockBooking, l *mocklogger.MockLogger) {
				s.EXPECT().
					DeleteBooking(gomock.Any(), 1).
					Return(errors.New("something went wrong"))
				l.EXPECT().
					Error("error", "something went wrong")
			},
			requestMethod:        "DELETE",
			expectedCodeStatus:   400,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLogger := mocklogger.NewMockLogger(ctrl)
			mockBookingService := mocksbookingservice.NewMockBooking(ctrl)

			test.mockBehavior(mockBookingService, mockLogger)

			bookingService := &services.Service{Booking: mockBookingService}
			handler := &Handler{service: bookingService, logger: mockLogger}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.requestMethod, "/bookings/1", nil)

			router := gin.New()
			router.Handle(test.requestMethod, "/bookings/:id", handler.deleteBooking)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCodeStatus, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())

		})
	}
}
