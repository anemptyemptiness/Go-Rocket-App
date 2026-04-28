package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1/mocks"
	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.GetOrderParams
	}

	type expected struct {
		resp orderv1.GetOrderRes
		err  error
	}

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		hullUUID        = uuid.New()
		engineUUID      = uuid.New()
		shieldUUID      = uuid.New()
		weaponUUID      = uuid.New()
		transactionUUID = uuid.New()
		paymentMethod   = model.PaymentMethodStringCard
		now             = time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
		unexpectedErr   = errors.New("внезапность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.OrderService)
	}{
		{
			name: "успешное получение заказа",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: &orderv1.OrderDto{
					OrderUUID:       orderUUID,
					HullUUID:        hullUUID,
					EngineUUID:      engineUUID,
					ShieldUUID:      orderv1.NewOptNilUUID(shieldUUID),
					WeaponUUID:      orderv1.NewOptNilUUID(weaponUUID),
					TotalPrice:      777000,
					TransactionUUID: orderv1.NewOptNilUUID(transactionUUID),
					PaymentMethod:   orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(paymentMethod)),
					Status:          orderv1.OrderStatusPAID,
					CreatedAt:       now,
				},
				err: nil,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID:       orderUUID,
						HullUUID:        hullUUID,
						EngineUUID:      engineUUID,
						ShieldUUID:      &shieldUUID,
						WeaponUUID:      &weaponUUID,
						TotalPrice:      777000,
						TransactionUUID: &transactionUUID,
						PaymentMethod:   &paymentMethod,
						Status:          model.OrderStatusPaid,
						CreatedAt:       now,
					}, nil)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrOrderNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID).
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

			resp, err := api.GetOrder(ctx, tc.args.params)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, &orderv1.OrderDto{}, resp)

				dto, ok := resp.(*orderv1.OrderDto)
				require.True(t, ok)
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetOrderUUID(), dto.GetOrderUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetHullUUID(), dto.GetHullUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetEngineUUID(), dto.GetEngineUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetTotalPrice(), dto.GetTotalPrice())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetStatus(), dto.GetStatus())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetCreatedAt(), dto.GetCreatedAt())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetShieldUUID(), dto.GetShieldUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetWeaponUUID(), dto.GetWeaponUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetTransactionUUID(), dto.GetTransactionUUID())
				assert.Equal(t, tc.expected.resp.(*orderv1.OrderDto).GetPaymentMethod(), dto.GetPaymentMethod())
			}
		})
	}
}
