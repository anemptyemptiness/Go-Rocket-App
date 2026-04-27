package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
		method    model.PaymentMethodString
	}

	type expected struct {
		transactionUUID uuid.UUID
		wantErr         error
	}

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		hullUUID        = uuid.New()
		engineUUID      = uuid.New()
		transactionUUID = uuid.New()
		unexpectedErr   = errors.New("внезапность")
		now             = time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient)
	}{
		{
			name: "успешная оплата заказа",
			args: args{
				orderUUID: orderUUID,
				method:    model.PaymentMethodStringCard,
			},
			expected: expected{
				transactionUUID: transactionUUID,
				wantErr:         nil,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, _ *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID:  orderUUID,
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						Status:     model.OrderStatusPendingPayment,
						CreatedAt:  now,
					}, nil)

				payment.EXPECT().
					PayOrder(mock.Anything, orderUUID, model.PaymentMethodStringCard).
					Return(transactionUUID, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.Status == model.OrderStatusPaid &&
							order.TransactionUUID != nil &&
							*order.TransactionUUID == transactionUUID &&
							order.PaymentMethod != nil &&
							*order.PaymentMethod == model.PaymentMethodStringCard
					})).
					Return(nil)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				orderUUID: orderUUID,
				method:    model.PaymentMethodStringCard,
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				wantErr:         ordererrs.ErrOrderNotFound,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: некорректный статус заказа",
			args: args{
				orderUUID: orderUUID,
				method:    model.PaymentMethodStringSBP,
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				wantErr:         ordererrs.ErrOrderStatusIncorrect,
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
			name: "ошибка: payment client invalid argument",
			args: args{
				orderUUID: orderUUID,
				method:    model.PaymentMethodStringCreditCard,
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				wantErr:         ordererrs.ErrPaymentClientInvalidArgument,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusPendingPayment,
					}, nil)

				payment.EXPECT().
					PayOrder(mock.Anything, orderUUID, model.PaymentMethodStringCreditCard).
					Return(uuid.Nil, ordererrs.ErrPaymentClientInvalidArgument)
			},
		},
		{
			name: "ошибка: обновить заказ в репозитории",
			args: args{
				orderUUID: orderUUID,
				method:    model.PaymentMethodStringInvestorMoney,
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				wantErr:         unexpectedErr,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						OrderUUID: orderUUID,
						Status:    model.OrderStatusPendingPayment,
					}, nil)

				payment.EXPECT().
					PayOrder(mock.Anything, orderUUID, model.PaymentMethodStringInvestorMoney).
					Return(transactionUUID, nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID == orderUUID &&
							order.Status == model.OrderStatusPaid &&
							order.TransactionUUID != nil &&
							*order.TransactionUUID == transactionUUID &&
							order.PaymentMethod != nil &&
							*order.PaymentMethod == model.PaymentMethodStringInvestorMoney
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

			resp, err := svc.Pay(ctx, tc.args.orderUUID, tc.args.method)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Equal(t, uuid.Nil, resp)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, resp)
				assert.Equal(t, tc.expected.transactionUUID, resp)
			}
		})
	}
}
