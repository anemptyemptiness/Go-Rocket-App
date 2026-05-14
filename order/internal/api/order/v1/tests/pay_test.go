package tests

import (
	"context"
	"net/http"
	"testing"

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

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req    *orderv1.PayOrderRequest
		params orderv1.PayOrderParams
	}

	type expected struct {
		res     orderv1.PayOrderRes
		wantErr bool
	}

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
		unexpectedErr   = assert.AnError
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
				res: &orderv1.PayOrderResponse{
					TransactionUUID: transactionUUID,
				},
				wantErr: false,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return(transactionUUID.String(), nil)
			},
		},
		{
			name: "ошибка: bad request конвертер",
			args: args{
				req: nil,
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.PayOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: ordererrs.ErrEmptyRequest.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {},
		},
		{
			name: "ошибка: bad request сервиса",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodSBP,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.PayOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: ordererrs.ErrOrderStatusIncorrect.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodSBP).
					Return("", pkgerr.InvalidArgument(ordererrs.ErrOrderStatusIncorrect))
			},
		},
		{
			name: "ошибка: not found",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCREDITCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.PayOrderNotFound{
					Code:    http.StatusNotFound,
					Message: ordererrs.ErrOrderNotFound.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCreditCard).
					Return("", pkgerr.NotFound(ordererrs.ErrOrderNotFound))
			},
		},
		{
			name: "ошибка: internal",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.PayOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: unexpectedErr.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return("", pkgerr.Internal(unexpectedErr))
			},
		},
		{
			name: "ошибка: conflict",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				res: &orderv1.PayOrderConflict{
					Code:    http.StatusConflict,
					Message: unexpectedErr.Error(),
				},
				wantErr: true,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID.String(), model.PaymentMethodCard).
					Return("", pkgerr.Conflict(unexpectedErr))
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
			require.NoError(t, err)
			require.Equal(t, tc.expected.res, resp)

			if !tc.expected.wantErr {
				result, ok := resp.(*orderv1.PayOrderResponse)
				require.True(t, ok)
				assert.Equal(t, result.TransactionUUID.String(), transactionUUID.String())
			}
		})
	}
}
