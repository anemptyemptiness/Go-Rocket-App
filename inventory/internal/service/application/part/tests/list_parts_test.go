package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	inventorysvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

func TestListParts(t *testing.T) {
	t.Parallel()

	type args struct {
		filter input.PartFilter
	}

	type expected struct {
		parts   []entity.Part
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

	hull := entity.RestorePart(hullUUID, "hull", "hull-desc", 10000, 10, 0, valueobject.PartTypeHull, nil, now)
	engine := entity.RestorePart(engineUUID, "engine", "engine-desc", 20000, 20, 0, valueobject.PartTypeEngine, nil, now)
	shield := entity.RestorePart(shieldUUID, "shield", "shield-desc", 30000, 30, 0, valueobject.PartTypeShield, nil, now)
	weapon := entity.RestorePart(weaponUUID, "weapon", "weapon-desc", 40000, 40, 0, valueobject.PartTypeWeapon, nil, now)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.Repository)
	}{
		{
			name: "успешное получение списка деталей",
			args: args{
				filter: input.PartFilter{
					UUIDs:    []string{hullUUID, engineUUID, shieldUUID, weaponUUID},
					PartType: valueobject.PartTypeUnspecified,
				},
			},
			expected: expected{
				parts: []entity.Part{
					hull,
					engine,
					shield,
					weapon,
				},
				wantErr: nil,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs:    []string{hullUUID, engineUUID, shieldUUID, weaponUUID},
						PartType: valueobject.PartTypeUnspecified,
					}).
					Return([]entity.Part{
						hull,
						engine,
						shield,
						weapon,
					}, nil)
			},
		},
		{
			name: "ошибка: репозиторий не нашел деталь",
			args: args{
				filter: input.PartFilter{
					UUIDs:    []string{shieldUUID, hullUUID},
					PartType: valueobject.PartTypeUnspecified,
				},
			},
			expected: expected{
				parts:   nil,
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs:    []string{shieldUUID, hullUUID},
						PartType: valueobject.PartTypeUnspecified,
					}).
					Return(nil, errs.ErrPartNotFound)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			compChecker := mocks.NewCompatibilityChecker(t)
			txManager := mocks.NewTxManager(t)

			tc.setupMock(repo)

			svc := inventorysvc.New(repo, compChecker, txManager)

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
				assert.IsType(t, []entity.Part{}, resp)
				require.Len(t, resp, len(tc.expected.parts))

				for idx, part := range resp {
					require.Less(t, idx, len(tc.expected.parts))
					require.NotEmpty(t, part)
					assert.Equal(t, tc.expected.parts[idx].GetPartUUID(), part.GetPartUUID())
				}
			}
		})
	}
}
