package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	inventorysvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/part"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/part/mocks"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	type args struct {
		uuidStr string
	}

	type expected struct {
		part    model.Part
		wantErr error
	}

	var (
		ctx           = context.Background()
		hullUUID      = gofakeit.UUID()
		hullUUIDEmpty = ""
		now           = time.Now()
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.Repository)
	}{
		{
			name: "успешное получение детали",
			args: args{
				uuidStr: hullUUID,
			},
			expected: expected{
				part: model.Part{
					UUID:          hullUUID,
					Name:          "hull",
					Description:   "hull-desc",
					Price:         10000,
					PartType:      model.PartTypeHull,
					StockQuantity: 10,
					CreatedAt:     now,
				},
				wantErr: nil,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					GetPart(ctx, hullUUID).
					Return(model.Part{
						UUID:          hullUUID,
						Name:          "hull",
						Description:   "hull-desc",
						Price:         10000,
						PartType:      model.PartTypeHull,
						StockQuantity: 10,
						CreatedAt:     now,
					}, nil)
			},
		},
		{
			name: "ошибка: UUID детали пустой",
			args: args{
				uuidStr: hullUUIDEmpty,
			},
			expected: expected{
				part:    model.Part{},
				wantErr: errs.ErrPartUUIDIsEmpty,
			},
			setupMock: func(repo *mocks.Repository) {},
		},
		{
			name: "ошибка: деталь не найдена в репозитории",
			args: args{
				uuidStr: hullUUID,
			},
			expected: expected{
				part:    model.Part{},
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					GetPart(ctx, hullUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			tc.setupMock(repo)

			svc := inventorysvc.New(repo)

			resp, err := svc.GetPart(ctx, tc.args.uuidStr)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, model.Part{}, resp)
				assert.Equal(t, resp.UUID, tc.expected.part.UUID)
			}
		})
	}
}
