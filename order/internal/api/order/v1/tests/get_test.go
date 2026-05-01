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
		resp *orderv1.OrderDto
		err  error
	}

	var (
		ctx                = context.Background()
		orderUUID          = uuid.New()
		hullUUID           = uuid.New()
		engineUUID         = uuid.New()
		shieldUUID         = uuid.New()
		weaponUUID         = uuid.New()
		transactionUUID    = uuid.New()
		transactionUUIDStr = transactionUUID.String()
		paymentMethod      = model.PaymentMethodCard
		now                = time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
		unexpectedErr      = errors.New("внезапность")
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
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
					Return(model.Order{
						UUID: orderUUID.String(),
						Items: []model.OrderItem{
							{PartUuid: hullUUID.String(), PartType: model.PartTypeHull, Price: 100000},
							{PartUuid: engineUUID.String(), PartType: model.PartTypeEngine, Price: 200000},
							{PartUuid: shieldUUID.String(), PartType: model.PartTypeShield, Price: 300000},
							{PartUuid: weaponUUID.String(), PartType: model.PartTypeWeapon, Price: 177000},
						},
						TotalPrice:      777000,
						TransactionUUID: &transactionUUIDStr,
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
				err: ordererrs.ErrOrderNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
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
				err: unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
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
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			dto, ok := resp.(*orderv1.OrderDto)
			require.True(t, ok)
			assert.Equal(t, tc.expected.resp.OrderUUID, dto.OrderUUID)
			assert.Equal(t, tc.expected.resp.HullUUID, dto.HullUUID)
			assert.Equal(t, tc.expected.resp.EngineUUID, dto.EngineUUID)
			assert.Equal(t, tc.expected.resp.ShieldUUID, dto.ShieldUUID)
			assert.Equal(t, tc.expected.resp.WeaponUUID, dto.WeaponUUID)
			assert.Equal(t, tc.expected.resp.TotalPrice, dto.TotalPrice)
			assert.Equal(t, tc.expected.resp.TransactionUUID, dto.TransactionUUID)
			assert.Equal(t, tc.expected.resp.PaymentMethod, dto.PaymentMethod)
			assert.Equal(t, tc.expected.resp.Status, dto.Status)
			assert.Equal(t, tc.expected.resp.CreatedAt, dto.CreatedAt)
		})
	}
}
