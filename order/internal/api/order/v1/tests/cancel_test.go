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
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.CancelOrderParams
	}

	type expected struct {
		res orderv1.CancelOrderRes
	}

	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		unexpectedErr = assert.AnError
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
				res: &orderv1.CancelOrderResponse{},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID.String()).
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
				res: &orderv1.CancelOrderNotFound{
					Code:    http.StatusNotFound,
					Message: ordererrs.ErrOrderNotFound.Error(),
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID.String()).
					Return(pkgerr.NotFound(ordererrs.ErrOrderNotFound))
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
				res: &orderv1.CancelOrderConflict{
					Code:    http.StatusConflict,
					Message: ordererrs.ErrOrderAlreadyPaid.Error(),
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID.String()).
					Return(pkgerr.Conflict(ordererrs.ErrOrderAlreadyPaid))
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
				res: &orderv1.CancelOrderConflict{
					Code:    http.StatusConflict,
					Message: ordererrs.ErrOrderAlreadyCancelled.Error(),
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID.String()).
					Return(pkgerr.Conflict(ordererrs.ErrOrderAlreadyCancelled))
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
				res: &orderv1.CancelOrderInternalServerError{
					Code:    http.StatusInternalServerError,
					Message: unexpectedErr.Error(),
				},
			},
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Cancel(ctx, orderUUID.String()).
					Return(pkgerr.Internal(unexpectedErr))
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
			require.NoError(t, err)
			require.Equal(t, tc.expected.res, resp)
		})
	}
}
