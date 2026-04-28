package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		order   model.Order
		wantErr error
	}

	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		hullUUID      = uuid.New()
		engineUUID    = uuid.New()
		unexpectedErr = errors.New("внезапность")
		now           = time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository)
	}{
		{
			name: "успешное получение заказа",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				order: model.Order{
					OrderUUID:  orderUUID,
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					TotalPrice: 800000,
					Status:     model.OrderStatusPendingPayment,
					CreatedAt:  now,
				},
				wantErr: nil,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID:  orderUUID,
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						TotalPrice: 800000,
						Status:     model.OrderStatusPendingPayment,
						CreatedAt:  now,
					}, nil)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				order:   model.Order{},
				wantErr: ordererrs.ErrOrderNotFound,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: внутренняя ошибка репозитория",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				order:   model.Order{},
				wantErr: unexpectedErr,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewOrderRepository(t)
			tc.setupMock(repo)

			paymentClient := mocks.NewPaymentClient(t)
			inventoryClient := mocks.NewInventoryClient(t)
			svc := ordersvc.New(repo, paymentClient, inventoryClient)

			resp, err := svc.Get(ctx, tc.args.orderUUID)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, model.Order{}, resp)
				assert.Equal(t, tc.expected.order, resp)
			}
		})
	}
}
