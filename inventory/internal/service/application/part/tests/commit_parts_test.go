package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	inventorysvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

func Test_CommitParts(t *testing.T) {
	t.Parallel()

	type args struct {
		req input.CommitPartsRequest
	}

	type expected struct {
		err error
	}

	var (
		ctx                    = context.Background()
		hullUUID               = gofakeit.UUID()
		hullUUIDZeroStock      = gofakeit.UUID()
		engineUUID             = gofakeit.UUID()
		engineUUIDZeroReserved = gofakeit.UUID()
		unexpectedErr          = assert.AnError
		now                    = time.Now()
	)

	hull := entity.RestorePart(hullUUID, "hull", "hull-desc", 10000, 10, 1, valueobject.PartTypeHull, nil, now)
	hullZeroStock := entity.RestorePart(hullUUIDZeroStock, "hull", "hull-desc", 10000, 0, 1, valueobject.PartTypeHull, nil, now)
	engine := entity.RestorePart(engineUUID, "engine", "engine-desc", 20000, 20, 1, valueobject.PartTypeEngine, nil, now)
	engineZeroReserved := entity.RestorePart(engineUUIDZeroReserved, "engine", "engine-desc", 20000, 20, 0, valueobject.PartTypeEngine, nil, now)

	tests := []struct {
		name       string
		args       args
		expected   expected
		setupMocks func(txManager *mocks.TxManager, repo *mocks.Repository)
	}{
		{
			name: "успешное списание деталей",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: nil,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, nil)

				repo.EXPECT().
					CommitParts(ctx, input.CommitPartsRequest{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return(nil)
			},
		},
		{
			name: "ошибка: деталь не найдена (ListParts)",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: errs.ErrPartNotFound,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, errs.ErrPartNotFound)
			},
		},
		{
			name: "ошибка: неверные свойства детали (ListParts)",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: errs.ErrInvalidProperties,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, errs.ErrInvalidProperties)
			},
		},
		{
			name: "ошибка: неожиданная ошибка (ListParts)",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, unexpectedErr)
			},
		},
		{
			name: "ошибка: stock quantity детали равен нулю",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUIDZeroStock,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: errs.ErrStockQuantityOrReservedIsEmpty,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUIDZeroStock, engineUUID},
					}).Return([]entity.Part{hullZeroStock, engine}, nil)
			},
		},
		{
			name: "ошибка: reserved детали равен нулю",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUIDZeroReserved,
					},
				},
			},
			expected: expected{
				err: errs.ErrStockQuantityOrReservedIsEmpty,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUIDZeroReserved},
					}).Return([]entity.Part{hull, engineZeroReserved}, nil)
			},
		},
		{
			name: "ошибка: деталь не найдена (CommitParts)",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: errs.ErrPartNotFound,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, nil)

				repo.EXPECT().
					CommitParts(ctx, input.CommitPartsRequest{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return(errs.ErrPartNotFound)
			},
		},
		{
			name: "ошибка: неожиданная ошибка (CommitParts)",
			args: args{
				req: input.CommitPartsRequest{
					UUIDs: []string{
						hullUUID,
						engineUUID,
					},
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMocks: func(txManager *mocks.TxManager, repo *mocks.Repository) {
				txManager.EXPECT().
					Do(ctx, mock.Anything).
					RunAndReturn(func(
						txCtx context.Context, fn func(ctx context.Context) error,
					) error {
						return fn(ctx)
					})

				repo.EXPECT().
					ListParts(ctx, input.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return([]entity.Part{hull, engine}, nil)

				repo.EXPECT().
					CommitParts(ctx, input.CommitPartsRequest{
						UUIDs: []string{hullUUID, engineUUID},
					}).Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			txManager := mocks.NewTxManager(t)
			compChecker := mocks.NewCompatibilityChecker(t)
			repo := mocks.NewRepository(t)

			tc.setupMocks(txManager, repo)

			svc := inventorysvc.New(repo, compChecker, txManager)

			err := svc.CommitParts(ctx, tc.args.req)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
