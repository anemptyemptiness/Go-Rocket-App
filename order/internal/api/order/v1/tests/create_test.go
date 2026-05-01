package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1/mocks"
	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req *orderv1.CreateOrderRequest
	}

	type expected struct {
		resp *orderv1.CreateOrderResponse
		err  error
	}

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		hullUUID      = gofakeit.UUID()
		engineUUID    = gofakeit.UUID()
		shieldUUID    = gofakeit.UUID()
		weaponUUID    = gofakeit.UUID()
		unexpectedErr = errors.New("неожиданность")
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
				resp: &orderv1.CreateOrderResponse{
					OrderUUID:  uuid.MustParse(orderUUID),
					TotalPrice: 999000,
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
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
			name: "ошибка: пустой реквест",
			args: args{
				req: nil,
			},
			expected: expected{
				err: ordererrs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: обязательные параметры не переданы",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.Nil,
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				err: ordererrs.ErrHullUUIDAndEngineUUIDAreRequired,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   uuid.Nil.String(),
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, ordererrs.ErrHullUUIDAndEngineUUIDAreRequired)
			},
		},
		{
			name: "ошибка: деталь закончилась",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				err: ordererrs.ErrPartIsOver,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, ordererrs.ErrPartIsOver)
			},
		},
		{
			name: "ошибка: inventory client not found",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				err: ordererrs.ErrInventoryClientNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, ordererrs.ErrInventoryClientNotFound)
			},
		},
		{
			name: "ошибка: inventory client invalid argument",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				err: ordererrs.ErrInventoryClientInvalidArgument,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, ordererrs.ErrInventoryClientInvalidArgument)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.MustParse(hullUUID),
					EngineUUID: uuid.MustParse(engineUUID),
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
					}).
					Return(model.Order{}, unexpectedErr)
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
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			result, ok := resp.(*orderv1.CreateOrderResponse)
			require.True(t, ok)
			assert.Equal(t, tc.expected.resp.OrderUUID, result.OrderUUID)
			assert.Equal(t, tc.expected.resp.TotalPrice, result.TotalPrice)
		})
	}
}
