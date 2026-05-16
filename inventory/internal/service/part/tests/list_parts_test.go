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

func TestListParts(t *testing.T) {
	t.Parallel()

	type args struct {
		filter model.PartFilter
	}

	type expected struct {
		parts   []model.Part
		wantErr error
	}

	var (
		ctx = context.Background()

		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
		shieldUUID = gofakeit.UUID()
		weaponUUID = gofakeit.UUID()

		now = time.Now()
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.Repository)
	}{
		{
			name: "успешное получение списка деталей",
			args: args{
				filter: model.PartFilter{
					UUIDs:    []string{hullUUID, engineUUID, shieldUUID, weaponUUID},
					PartType: model.PartTypeUnspecified,
				},
			},
			expected: expected{
				parts: []model.Part{
					{
						UUID:          hullUUID,
						Name:          "hull",
						Description:   "hull-desc",
						Price:         10000,
						PartType:      model.PartTypeHull,
						StockQuantity: 10,
						CreatedAt:     now,
					},
					{
						UUID:          engineUUID,
						Name:          "engine",
						Description:   "engine-desc",
						Price:         20000,
						PartType:      model.PartTypeEngine,
						StockQuantity: 20,
						CreatedAt:     now,
					},
					{
						UUID:          shieldUUID,
						Name:          "shield",
						Description:   "shield-desc",
						Price:         30000,
						PartType:      model.PartTypeShield,
						StockQuantity: 30,
						CreatedAt:     now,
					},
					{
						UUID:          weaponUUID,
						Name:          "weapon",
						Description:   "weapon-desc",
						Price:         40000,
						PartType:      model.PartTypeWeapon,
						StockQuantity: 40,
						CreatedAt:     now,
					},
				},
				wantErr: nil,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					ListParts(ctx, model.PartFilter{
						UUIDs:    []string{hullUUID, engineUUID, shieldUUID, weaponUUID},
						PartType: model.PartTypeUnspecified,
					}).
					Return([]model.Part{
						{
							UUID:          hullUUID,
							Name:          "hull",
							Description:   "hull-desc",
							Price:         10000,
							PartType:      model.PartTypeHull,
							StockQuantity: 10,
							CreatedAt:     now,
						},
						{
							UUID:          engineUUID,
							Name:          "engine",
							Description:   "engine-desc",
							Price:         20000,
							PartType:      model.PartTypeEngine,
							StockQuantity: 20,
							CreatedAt:     now,
						},
						{
							UUID:          shieldUUID,
							Name:          "shield",
							Description:   "shield-desc",
							Price:         30000,
							PartType:      model.PartTypeShield,
							StockQuantity: 30,
							CreatedAt:     now,
						},
						{
							UUID:          weaponUUID,
							Name:          "weapon",
							Description:   "weapon-desc",
							Price:         40000,
							PartType:      model.PartTypeWeapon,
							StockQuantity: 40,
							CreatedAt:     now,
						},
					}, nil)
			},
		},
		{
			name: "ошибка: репозиторий не нашел деталь",
			args: args{
				filter: model.PartFilter{
					UUIDs:    []string{shieldUUID, hullUUID},
					PartType: model.PartTypeUnspecified,
				},
			},
			expected: expected{
				parts:   nil,
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					ListParts(ctx, model.PartFilter{
						UUIDs:    []string{shieldUUID, hullUUID},
						PartType: model.PartTypeUnspecified,
					}).
					Return(nil, errs.ErrPartNotFound)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			tc.setupMock(repo)

			svc := inventorysvc.New(repo)

			resp, err := svc.ListParts(ctx, tc.args.filter)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp)
				assert.IsType(t, []model.Part{}, resp)
				require.Len(t, resp, len(tc.expected.parts))

				for idx, part := range resp {
					require.Less(t, idx, len(tc.expected.parts))
					require.NotEmpty(t, part)
					assert.Equal(t, tc.expected.parts[idx].UUID, part.UUID)
				}
			}
		})
	}
}
