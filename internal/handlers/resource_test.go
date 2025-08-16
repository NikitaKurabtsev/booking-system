package handlers

import (
	"errors"
	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/NikitaKurabtsev/booking-system/internal/services"
	resourceservicemock "github.com/NikitaKurabtsev/booking-system/internal/services/mocks"
	loggermock "github.com/NikitaKurabtsev/booking-system/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_getAllResources(t *testing.T) {
	type mockBehavior func(s *resourceservicemock.MockResource, l *loggermock.MockLogger)

	mockData := []domain.Resource{
		{ID: 1, Name: "office", Type: "available"},
	}

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedCodeStatus   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *resourceservicemock.MockResource, l *loggermock.MockLogger) {
				s.EXPECT().GetResources(gomock.Any()).Return(mockData, nil)
			},
			expectedCodeStatus:   200,
			expectedResponseBody: `{"data":[{"id":1,"name":"office","type":"available"}]}`,
		},
		{
			name: "Service Error",
			mockBehavior: func(s *resourceservicemock.MockResource, l *loggermock.MockLogger) {
				s.EXPECT().GetResources(gomock.Any()).Return(nil, errors.New("something went wrong"))
				l.EXPECT().Error("error", "something went wrong")
			},
			expectedCodeStatus:   500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLogger := loggermock.NewMockLogger(ctrl)
			mockResourceService := resourceservicemock.NewMockResource(ctrl)

			test.mockBehavior(mockResourceService, mockLogger)

			service := &services.Service{Resource: mockResourceService}

			handler := &Handler{
				logger:  mockLogger,
				service: service,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/resources", nil)

			router := gin.New()
			router.GET("/resources", handler.getAllResources)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedResponseBody, w.Body.String())
			assert.Equal(t, test.expectedCodeStatus, w.Code)
		})
	}
}
