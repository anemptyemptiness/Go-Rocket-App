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

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req    *orderv1.PayOrderRequest
		params orderv1.PayOrderParams
	}

	type expected struct {
		resp *orderv1.PayOrderResponse
		err  error
	}

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
		unexpectedErr   = errors.New("внезапность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.OrderService)
	}{
		{
			name: "успешная оплата заказа",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: &orderv1.PayOrderResponse{
					TransactionUUID: transactionUUID,
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return(transactionUUID.String(), nil)
			},
		},
		{
			name: "ошибка: пустой реквест",
			args: args{
				req: nil,
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				err: ordererrs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: некорректный статус заказа",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodSBP,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				err: ordererrs.ErrOrderStatusIncorrect,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodSBP).
					Return("", ordererrs.ErrOrderStatusIncorrect)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCREDITCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				err: ordererrs.ErrOrderNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCreditCard).
					Return("", ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: payment client invalid argument",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodINVESTORMONEY,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				err: ordererrs.ErrPaymentClientInvalidArgument,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodInvestorMoney).
					Return("", ordererrs.ErrPaymentClientInvalidArgument)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				err: unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return("", unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.NewAPI(svc)

			resp, err := api.PayOrder(ctx, tc.args.req, tc.args.params)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			result, ok := resp.(*orderv1.PayOrderResponse)
			require.True(t, ok)
			assert.Equal(t, tc.expected.resp.TransactionUUID, result.TransactionUUID)
		})
	}
}
