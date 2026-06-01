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
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	type args struct {
		uuidStr string
	}

	type expected struct {
		part    entity.Part
		wantErr error
	}

	var (
		ctx           = context.Background()
		hullUUID      = gofakeit.UUID()
		hullUUIDEmpty = ""
		now           = time.Now()
	)

	part := entity.RestorePart(
		hullUUID,
		"hull",
		"hull-desc",
		10000,
		10,
		0,
		valueobject.PartTypeHull,
		nil,
		now,
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
				part:    part,
				wantErr: nil,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					GetPart(ctx, hullUUID).
					Return(part, nil)
			},
		},
		{
			name: "ошибка: UUID детали пустой",
			args: args{
				uuidStr: hullUUIDEmpty,
			},
			expected: expected{
				part:    entity.Part{},
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
				part:    entity.Part{},
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(repo *mocks.Repository) {
				repo.EXPECT().
					GetPart(ctx, hullUUID).
					Return(entity.Part{}, errs.ErrPartNotFound)
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

			resp, err := svc.GetPart(ctx, tc.args.uuidStr)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, entity.Part{}, resp)
				assert.Equal(t, resp.GetPartUUID(), tc.expected.part.GetPartUUID())
			}
		})
	}
}
