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
		resp orderv1.PayOrderRes
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
				err: nil,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodStringCard).
					Return(transactionUUID, nil)
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
				resp: nil,
				err:  ordererrs.ErrEmptyRequest,
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
				resp: nil,
				err:  ordererrs.ErrOrderStatusIncorrect,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodStringSBP).
					Return(uuid.Nil, ordererrs.ErrOrderStatusIncorrect)
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
				resp: nil,
				err:  ordererrs.ErrOrderNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodStringCreditCard).
					Return(uuid.Nil, ordererrs.ErrOrderNotFound)
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
				resp: nil,
				err:  ordererrs.ErrPaymentClientInvalidArgument,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodStringInvestorMoney).
					Return(uuid.Nil, ordererrs.ErrPaymentClientInvalidArgument)
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
				resp: nil,
				err:  unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, model.PaymentMethodStringCard).
					Return(uuid.Nil, unexpectedErr)
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
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, &orderv1.PayOrderResponse{}, resp)

				result, ok := resp.(*orderv1.PayOrderResponse)
				require.True(t, ok)
				assert.Equal(t, tc.expected.resp.(*orderv1.PayOrderResponse).GetTransactionUUID(), result.GetTransactionUUID())
			}
		})
	}
}
