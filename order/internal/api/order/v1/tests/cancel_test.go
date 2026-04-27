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
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.CancelOrderParams
	}

	type expected struct {
		resp orderv1.CancelOrderRes
		err  error
	}

	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		unexpectedErr = errors.New("неожиданность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.OrderService)
	}{
		{
			name: "успешная отмена заказа",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: &orderv1.CancelOrderResponse{},
				err:  nil,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(nil)
			},
		},
		{
			name: "ошибка: заказ не найден",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrOrderNotFound,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(ordererrs.ErrOrderNotFound)
			},
		},
		{
			name: "ошибка: заказ уже оплачен",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrOrderAlreadyPaid,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(ordererrs.ErrOrderAlreadyPaid)
			},
		},
		{
			name: "ошибка: заказ уже отменён",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  ordererrs.ErrOrderAlreadyCancelled,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(ordererrs.ErrOrderAlreadyCancelled)
			},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{
				resp: nil,
				err:  unexpectedErr,
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID).
					Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderapi.NewAPI(svc)

			resp, err := api.CancelOrder(ctx, tc.args.params)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.IsType(t, &orderv1.CancelOrderResponse{}, resp)
			}
		})
	}
}
