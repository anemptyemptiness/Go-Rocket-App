package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.GetPartRequest
	}

	type expected struct {
		resp    *inventoryv1.GetPartResponse
		wantErr error
	}

	var (
		ctx = context.Background()

		hullUUID        = gofakeit.UUID()
		hullUUIDEmpty   = ""
		hullUUIDInvalid = "flasmfkamsf"

		internalError = errors.New("внезапность")

		now = time.Now()
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.InventoryService)
	}{
		{
			name: "успешное получение детали",
			args: args{
				req: &inventoryv1.GetPartRequest{
					Uuid: hullUUID,
				},
			},
			expected: expected{
				resp: &inventoryv1.GetPartResponse{
					Part: &inventoryv1.Part{
						Uuid:          hullUUID,
						Name:          "hull",
						Description:   "hull-desc",
						Price:         10000,
						PartType:      inventoryv1.PartType_PART_TYPE_HULL,
						StockQuantity: 10,
						CreatedAt:     timestamppb.New(now),
					},
				},
				wantErr: nil,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
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
			name: "ошибка: пустой реквест",
			args: args{
				req: nil,
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.InventoryService) {},
		},
		{
			name: "ошибка: деталь не найдена",
			args: args{
				req: &inventoryv1.GetPartRequest{
					Uuid: hullUUID,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrPartNotFound,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					GetPart(ctx, hullUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
		},
		{
			name: "ошибка: UUID детали пустой",
			args: args{
				req: &inventoryv1.GetPartRequest{
					Uuid: hullUUIDEmpty,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrPartUUIDIsEmpty,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					GetPart(ctx, hullUUIDEmpty).
					Return(model.Part{}, errs.ErrPartUUIDIsEmpty)
			},
		},
		{
			name: "ошибка: UUID детали некорректный",
			args: args{
				req: &inventoryv1.GetPartRequest{
					Uuid: hullUUIDInvalid,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: errs.ErrIncorrectPartUUID,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					GetPart(ctx, hullUUIDInvalid).
					Return(model.Part{}, errs.ErrIncorrectPartUUID)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &inventoryv1.GetPartRequest{
					Uuid: hullUUID,
				},
			},
			expected: expected{
				resp:    nil,
				wantErr: internalError,
			},
			setupMock: func(svc *mocks.InventoryService) {
				svc.EXPECT().
					GetPart(ctx, hullUUID).
					Return(model.Part{}, internalError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewInventoryService(t)
			tc.setupMock(svc)

			api := inventoryapi.New(svc)

			resp, err := api.GetPart(ctx, tc.args.req)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.IsType(t, &inventoryv1.GetPartResponse{}, resp)
				assert.NotEmpty(t, resp)
				assert.NotNil(t, resp.GetPart())
				assert.Equal(t, tc.expected.resp.GetPart().GetUuid(), resp.GetPart().GetUuid())
			}
		})
	}
}
