package tests

import (
	"context"
	"errors"
	"testing"

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
		resp orderv1.CreateOrderRes
		err  error
	}

	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		hullUUID      = uuid.New()
		engineUUID    = uuid.New()
		shieldUUID    = uuid.New()
		weaponUUID    = uuid.New()
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
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					ShieldUUID: orderv1.NewOptNilUUID(shieldUUID),
					WeaponUUID: orderv1.NewOptNilUUID(weaponUUID),
				},
			},
			expected: expected{
				resp: &orderv1.CreateOrderResponse{
					OrderUUID:  orderUUID,
					TotalPrice: 999000,
				},
				err: nil,
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
						OrderUUID:  orderUUID,
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
				resp: nil,
				err:  ordererrs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: обязательные параметры не переданы",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   uuid.Nil,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrHullUUIDAndEngineUUIDAreRequired,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   uuid.Nil,
						EngineUUID: engineUUID,
						ShieldUUID: nil,
						WeaponUUID: nil,
					}).
					Return(model.Order{}, ordererrs.ErrHullUUIDAndEngineUUIDAreRequired)
			},
		},
		{
			name: "ошибка: деталь закончилась",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrPartIsOver,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						ShieldUUID: nil,
						WeaponUUID: nil,
					}).
					Return(model.Order{}, ordererrs.ErrPartIsOver)
			},
		},
		{
			name: "ошибка: inventory client not found",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrInventoryClientNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						ShieldUUID: nil,
						WeaponUUID: nil,
					}).
					Return(model.Order{}, ordererrs.ErrInventoryClientNotFound)
			},
		},
		{
			name: "ошибка: inventory client invalid argument",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrInventoryClientInvalidArgument,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						ShieldUUID: nil,
						WeaponUUID: nil,
					}).
					Return(model.Order{}, ordererrs.ErrInventoryClientInvalidArgument)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, model.CreateOrderRequest{
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						ShieldUUID: nil,
						WeaponUUID: nil,
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
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, &orderv1.CreateOrderResponse{}, resp)

				result, ok := resp.(*orderv1.CreateOrderResponse)
				require.True(t, ok)
				assert.Equal(t, tc.expected.resp.(*orderv1.CreateOrderResponse).GetOrderUUID(), result.GetOrderUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.CreateOrderResponse).GetTotalPrice(), result.GetTotalPrice())
			}
		})
	}
}
