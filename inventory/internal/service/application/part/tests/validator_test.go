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
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func TestValidateCompatibility(t *testing.T) {
	t.Parallel()

	type args struct {
		req input.ValidateCompatibilityRequest
	}

	type expected struct {
		wantErr error
	}

	var (
		ctx        = context.Background()
		now        = time.Now()
		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
		shieldUUID = gofakeit.UUID()
		weaponUUID = gofakeit.UUID()
	)

	hull := entity.RestorePart(hullUUID, "name", "desc", 10, 10, 0, valueobject.PartTypeHull, nil, now)
	engine := entity.RestorePart(engineUUID, "name", "desc", 10, 10, 0, valueobject.PartTypeEngine, nil, now)
	shield := entity.RestorePart(shieldUUID, "name", "desc", 10, 10, 0, valueobject.PartTypeWeapon, nil, now)
	weapon := entity.RestorePart(weaponUUID, "name", "desc", 10, 10, 0, valueobject.PartTypeWeapon, nil, now)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.Repository, compChecker *mocks.CompatibilityChecker)
	}{
		{
			name: "успешная валидация",
			args: args{
				req: input.ValidateCompatibilityRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				wantErr: nil,
			},
			setupMock: func(repo *mocks.Repository, compChecker *mocks.CompatibilityChecker) {
				parts := []entity.Part{hull, engine}

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: input.ValidateCompatibilityRequest{
							HullUUID:   hullUUID,
							EngineUUID: engineUUID,
						}.UUIDs(),
					}).
					Return(parts, nil)

				compChecker.EXPECT().
					Check(parts).
					Return(nil)
			},
		},
		{
			name: "ошибка: выбраны несоответствующие типы деталей",
			args: args{
				req: input.ValidateCompatibilityRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					ShieldUUID: shieldUUID,
					WeaponUUID: weaponUUID,
				},
			},
			expected: expected{
				wantErr: pkgerr.InvalidArgument(errs.ErrPartTypeMismatch),
			},
			setupMock: func(repo *mocks.Repository, _ *mocks.CompatibilityChecker) {
				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: input.ValidateCompatibilityRequest{
							HullUUID:   hullUUID,
							EngineUUID: engineUUID,
							ShieldUUID: shieldUUID,
							WeaponUUID: weaponUUID,
						}.UUIDs(),
					}).
					Return([]entity.Part{hull, engine, shield, weapon}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewRepository(t)
			compChecker := mocks.NewCompatibilityChecker(t)
			txManager := mocks.NewTxManager(t)

			svc := inventorysvc.New(repo, compChecker, txManager)

			tc.setupMock(repo, compChecker)

			err := svc.ValidateCompatibility(ctx, tc.args.req)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expected.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
