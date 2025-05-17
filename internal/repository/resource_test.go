package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/NikitaKurabtsev/booking-system/internal/models"
	mock_cache "github.com/NikitaKurabtsev/booking-system/internal/repository/mocks"
	"github.com/NikitaKurabtsev/booking-system/pkg/db"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestResourceRepository_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock_cache.NewMockCache(ctrl)

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error occured when opening mock gpx pool: '%s'", err)
	}
	defer mockPool.Close()

	r := NewResourceRepository(mockPool, mockCache)

	tests := []struct {
		name    string
		mock    func()
		want    []models.Resource
		wantErr bool
	}{
		{
			name: "Ok_CacheMiss",
			mock: func() {
				mockCache.EXPECT().
					Get(gomock.Any(), "resources:all").
					Return("", errors.New("not found"))

				rawQuery := `
					SELECT id, name, type, status
					FROM %s
				`
				query := fmt.Sprintf(rawQuery, db.ResourcesTable)

				rows := pgxmock.NewRows([]string{"id", "name", "type", "status"}).
					AddRow(1, "testResource1", "office", "available").
					AddRow(2, "testResource2", "commercial", "available").
					AddRow(3, "testResource3", "office", "booked")

				mockPool.ExpectQuery(query).WillReturnRows(rows)

				mockCache.EXPECT().
					Set(gomock.Any(), "resources:all", gomock.Any(), gomock.Any()).
					AnyTimes()
			},
			want: []models.Resource{
				{ID: 1, Name: "testResource1", Type: "office", Status: "available"},
				{ID: 2, Name: "testResource2", Type: "commercial", Status: "available"},
				{ID: 3, Name: "testResource3", Type: "office", Status: "booked"},
			},
		},
		{
			name: "Ok_CacheHit",
			mock: func() {
				resources := []models.Resource{
					{ID: 1, Name: "testResource1", Type: "office", Status: "available"},
				}
				cached, _ := json.Marshal(resources)

				mockCache.EXPECT().
					Get(gomock.Any(), "resources:all").
					Return(string(cached), nil)

			},
			want: []models.Resource{
				{ID: 1, Name: "testResource1", Type: "office", Status: "available"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}
