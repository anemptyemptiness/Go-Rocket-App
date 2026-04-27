package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		wantErr error
	}

	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		unexpectedErr = errors.New("внезапность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient)
	}{
		{
			name: "успешная отмена заказа",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				wantErr: nil,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusPendingPayment,
					}, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.Status == model.OrderStatusCancelled
					})).
					Return(nil)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				wantErr: ordererrs.ErrOrderNotFound,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: заказ уже оплачен",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				wantErr: ordererrs.ErrOrderAlreadyPaid,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusPaid,
					}, nil)
			},
		},
		{
			name: "ошибка: заказ уже отменён",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				wantErr: ordererrs.ErrOrderAlreadyCancelled,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusCancelled,
					}, nil)
			},
		},
		{
			name: "ошибка: обновить заказ в репозитории",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{
				wantErr: unexpectedErr,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusPendingPayment,
					}, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.Status == model.OrderStatusCancelled
					})).
					Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewOrderRepository(t)
			paymentClient := mocks.NewPaymentClient(t)
			inventoryClient := mocks.NewInventoryClient(t)
			tc.setupMock(repo, paymentClient, inventoryClient)

			svc := ordersvc.New(repo, paymentClient, inventoryClient)

			err := svc.Cancel(ctx, tc.args.orderUUID)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
