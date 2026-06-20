package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1/mocks"
	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req *orderv1.CreateOrderRequest
	}

	type expected struct {
		res     orderv1.CreateOrderRes
		wantErr bool
	}

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		hullUUID      = gofakeit.UUID()
		engineUUID    = gofakeit.UUID()
		shieldUUID    = gofakeit.UUID()
		weaponUUID    = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.OrderService)
	}{
		{
			name: "успешное создание заказа",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
					ShieldUUID: orderv1.NewOptNilUUID(uuid.MustParse(shieldUUID)),
					WeaponUUID: orderv1.NewOptNilUUID(uuid.MustParse(weaponUUID)),
				},
			},
			expected: expected{
				res: &orderv1.CreateOrderResponse{
					OrderUUID:  uuid.MustParse(orderUUID),
					TotalPrice: 999000,
				},
				wantErr: false,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						ShieldUUID: &shieldUUID,
						WeaponUUID: &weaponUUID,
					}).
					Return(model.Order{
						UUID:       orderUUID,
						TotalPrice: 999000,
					}, nil)
			},
		},
		{
			name: "ошибка: bad request конвертера",
			args: args{
				req: nil,
			},
			expected: expected{
				res: &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: ordererrs.ErrEmptyRequest.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: bad request",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.Nil,
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				res: &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: ordererrs.ErrHullUUIDAndEngineUUIDAreRequired.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: conflict",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				res: &orderv1.CreateOrderConflict{
					Code:    http.StatusConflict,
					Message: ordererrs.ErrPartIsOver.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, pkgerr.Conflict(ordererrs.ErrPartIsOver))
			},
		},
		{
			name: "ошибка: internal",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				res: &orderv1.CreateOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: unexpectedErr.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, pkgerr.Internal(unexpectedErr))
			},
		},
		{
			name: "ошибка: not found",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				res: &orderv1.CreateOrderNotFound{
					Code:    http.StatusNotFound,
					Message: assert.AnError.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, pkgerr.NotFound(assert.AnError))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.NewAPI(svc)

			resp, err := api.CreateOrder(ctx, tc.args.req)
			require.Equal(t, tc.expected.res, resp)
			require.NoError(t, err)

			if !tc.expected.wantErr {
				result, ok := resp.(*orderv1.CreateOrderResponse)
				require.True(t, ok)
				assert.Equal(t, result.OrderUUID.String(), orderUUID)
				assert.Equal(t, result.TotalPrice, int64(999000))
			}
		})
	}
}
