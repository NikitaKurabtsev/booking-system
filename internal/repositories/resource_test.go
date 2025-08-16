package repositories

import (
	"context"
	"encoding/json"
	"errors"
	cacheMocks "github.com/NikitaKurabtsev/booking-system/pkg/cache/mocks"
	"testing"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestResourceRepository_GetAll(t *testing.T) {
	const (
		cacheKey = "resources:all"
		query    = `
			SELECT id, name, type
			FROM resources
`
	)

	var ErrNotFound = errors.New("not found")

	testResources := []domain.Resource{
		{ID: 1, Name: "testResource1", Type: "office"},
		{ID: 2, Name: "testResource2", Type: "commercial"},
		{ID: 3, Name: "testResource3", Type: "office"},
	}

	ctrl := gomock.NewController(t)

	mockCache := cacheMocks.NewMockCache(ctrl)

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error occured when opening mock gpx pool: '%s'", err)
	}
	defer mockPool.Close()

	r := NewResourceRepository(mockPool, mockCache)

	tests := []struct {
		name    string
		mock    func()
		want    []domain.Resource
		wantErr bool
	}{
		{
			name: "Ok_CacheMiss",
			mock: func() {
				mockCache.EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", ErrNotFound)

				rows := pgxmock.NewRows([]string{"id", "name", "type"}).
					AddRow(testResources[0].ID, testResources[0].Name, testResources[0].Type).
					AddRow(testResources[1].ID, testResources[1].Name, testResources[1].Type).
					AddRow(testResources[2].ID, testResources[2].Name, testResources[2].Type)

				mockPool.ExpectQuery(query).WillReturnRows(rows)

				mockCache.EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), gomock.Any()).
					AnyTimes()
			},
			want: testResources,
		},
		{
			name: "Ok_CacheHit",
			mock: func() {
				cached, _ := json.Marshal(testResources)

				mockCache.EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(string(cached), nil)

			},
			want: testResources,
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
