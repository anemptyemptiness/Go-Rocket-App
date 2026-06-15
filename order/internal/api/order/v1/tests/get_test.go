package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1/mocks"
	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.GetOrderParams
	}

	type expected struct {
		res     orderv1.GetOrderRes
		wantErr bool
	}

	var (
		ctx                = context.Background()
		orderUUID          = uuid.New()
		hullUUID           = uuid.New()
		engineUUID         = uuid.New()
		userUUID           = uuid.New()
		shieldUUID         = uuid.New()
		weaponUUID         = uuid.New()
		transactionUUID    = uuid.New()
		transactionUUIDStr = transactionUUID.String()
		paymentMethod      = model.PaymentMethodCard
		now                = time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
		unexpectedErr      = assert.AnError
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
				res: &orderv1.OrderDto{
					OrderUUID:       orderUUID,
					HullUUID:        hullUUID,
					UserUUID:        userUUID,
					EngineUUID:      engineUUID,
					ShieldUUID:      orderv1.NewOptNilUUID(shieldUUID),
					WeaponUUID:      orderv1.NewOptNilUUID(weaponUUID),
					TotalPrice:      777000,
					TransactionUUID: orderv1.NewOptNilUUID(transactionUUID),
					PaymentMethod:   orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(paymentMethod)),
					Status:          orderv1.OrderStatusPAID,
					CreatedAt:       now,
				},
				wantErr: false,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
					Return(model.Order{
						UUID:     orderUUID.String(),
						UserUUID: userUUID.String(),
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
			name: "ошибка: not found",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.GetOrderNotFound{
					Code:    http.StatusNotFound,
					Message: ordererrs.ErrOrderNotFound.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
					Return(model.Order{}, pkgerr.NotFound(ordererrs.ErrOrderNotFound))
			},
		},
		{
			name: "ошибка: internal",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.GetOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: unexpectedErr.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Get(ctx, orderUUID.String()).
					Return(model.Order{}, pkgerr.Internal(unexpectedErr))
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
			require.NoError(t, err)
			require.Equal(t, tc.expected.res, resp)

			if !tc.expected.wantErr {
				result, ok := resp.(*orderv1.OrderDto)
				require.True(t, ok)
				assert.Equal(t, result.OrderUUID.String(), orderUUID.String())
				assert.Equal(t, result.HullUUID.String(), hullUUID.String())
				assert.Equal(t, result.EngineUUID.String(), engineUUID.String())
				assert.Equal(t, result.ShieldUUID.Value.String(), shieldUUID.String())
				assert.Equal(t, result.WeaponUUID.Value.String(), weaponUUID.String())
				assert.Equal(t, result.TotalPrice, int64(777000))
				assert.Equal(t, result.TransactionUUID.Value.String(), transactionUUID.String())
				assert.Equal(t, result.PaymentMethod.Value, orderv1.PaymentMethodCARD)
				assert.Equal(t, result.Status, orderv1.OrderStatusPAID)
				assert.False(t, result.CreatedAt.IsZero())
			}
		})
	}
}
