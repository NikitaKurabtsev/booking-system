package repositories

import (
	"context"
	mockcache "github.com/NikitaKurabtsev/booking-system/pkg/cache/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"regexp"
	"testing"
	"time"

	"github.com/NikitaKurabtsev/booking-system/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestBookingRepository_Create(t *testing.T) {
	const (
		cacheKey = "resources:all"

		overlapQuery = `
			SELECT start_time, end_time
			FROM bookings
			WHERE resource_id = $1
			AND $2 < end_time
			AND $3 > start_time
`
		createQuery = `
			INSERT INTO bookings 
				(resource_id, user_id, start_time, end_time)
			VALUES 
				($1, $2, $3, $4)
			RETURNING id
`
	)

	timeNow := time.Now().Truncate(time.Second).UTC()
	endTime := timeNow.Add(2 * time.Hour).UTC()

	testBookingInput := domain.Booking{
		ID:         1,
		ResourceID: 123,
		UserID:     1,
		StartTime:  timeNow,
		EndTime:    endTime,
	}

	testBookingOutput := booking{
		ID:         1,
		ResourceID: 123,
		UserID:     1,
		StartTime:  pgtype.Timestamp{Time: timeNow, Valid: true},
		EndTime:    pgtype.Timestamp{Time: endTime, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mockcache.NewMockCache(ctrl)
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error occured when opening mock gpx pool: %v", err)
	}
	defer mockPool.Close()

	r := NewBookingRepository(mockPool, mockCache)

	tests := []struct {
		name  string
		input domain.Booking
		//output  booking
		want    int
		wantErr bool
		errText string
		mock    func()
	}{
		{
			name:  "Ok",
			input: testBookingInput,
			want:  1,
			mock: func() {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(regexp.QuoteMeta(overlapQuery)).
					WithArgs(
						testBookingInput.ResourceID,
						testBookingInput.StartTime,
						testBookingInput.EndTime).
					WillReturnError(pgx.ErrNoRows)

				row := pgxmock.NewRows([]string{"id"}).AddRow(1)
				mockPool.ExpectQuery(regexp.QuoteMeta(createQuery)).
					WithArgs(
						testBookingInput.ResourceID,
						testBookingInput.UserID,
						// TODO: change output to the input
						// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!
						testBookingOutput.StartTime,
						testBookingOutput.EndTime).
					WillReturnRows(row)

				mockPool.ExpectCommit()

				mockCache.EXPECT().Delete(gomock.Any(), cacheKey).Return(nil).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(context.Background(), tt.input)
			if tt.wantErr {
				// ...
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}
