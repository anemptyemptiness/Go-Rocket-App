package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-faster/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	paymentapi "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1/mocks"
	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req *paymentv1.PayOrderRequest
	}

	type expected struct {
		transactionUUID string
		err             error
	}

	paymentMethods := []paymentv1.PaymentMethod{
		paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	var (
		ctx              = context.Background()
		orderUUID        = gofakeit.UUID()
		transactionUUID  = gofakeit.UUID()
		emptyOrderUUID   = ""
		invalidOrderUUID = "kfdnmskjfnsd"
		unexpectedErr    = errors.New("внезапность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.PaymentService)
	}{
		{
			name: "успешная оплата заказа",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID,
					PaymentMethod: paymentMethods[1],
				},
			},
			expected: expected{
				transactionUUID: transactionUUID,
				err:             nil,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, orderUUID, model.PaymentMethod(paymentMethods[1])).
					Return(transactionUUID, nil)
			},
		},
		{
			name: "успешно возвращенный transactionUUID",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID,
					PaymentMethod: paymentMethods[2],
				},
			},
			expected: expected{
				transactionUUID: transactionUUID,
				err:             nil,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, orderUUID, model.PaymentMethod(paymentMethods[2])).
					Return(transactionUUID, nil)
			},
		},
		{
			name: "ошибка: идентификатор заказа пустой",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     emptyOrderUUID,
					PaymentMethod: paymentMethods[1],
				},
			},
			expected: expected{
				transactionUUID: "",
				err:             errs.ErrOrderUUIDIsEmpty,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, emptyOrderUUID, model.PaymentMethod(paymentMethods[1])).
					Return("", errs.ErrOrderUUIDIsEmpty)
			},
		},
		{
			name: "ошибка: невалидный orderUUID",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     invalidOrderUUID,
					PaymentMethod: paymentMethods[1],
				},
			},
			expected: expected{
				transactionUUID: "",
				err:             errs.ErrIncorrectOrderUUID,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, invalidOrderUUID, model.PaymentMethod(paymentMethods[1])).
					Return("", errs.ErrIncorrectOrderUUID)
			},
		},
		{
			name: "ошибка: метод оплаты неопределённый",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID,
					PaymentMethod: paymentMethods[0],
				},
			},
			expected: expected{
				transactionUUID: "",
				err:             errs.ErrPaymentMethodUnspecified,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, orderUUID, model.PaymentMethod(paymentMethods[0])).
					Return("", errs.ErrPaymentMethodUnspecified)
			},
		},
		{
			name: "ошибка: пустой реквест",
			args: args{
				req: nil,
			},
			expected: expected{
				transactionUUID: "",
				err:             errs.ErrEmptyRequest,
			},
			setupMock: func(svc *mocks.PaymentService) {},
		},
		{
			name: "ошибка: внутренняя ошибка",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID,
					PaymentMethod: paymentMethods[3],
				},
			},
			expected: expected{
				transactionUUID: "",
				err:             unexpectedErr,
			},
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					PayOrder(ctx, orderUUID, model.PaymentMethod(paymentMethods[3])).
					Return("", unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewPaymentService(t)
			tc.setupMock(svc)

			api := paymentapi.New(svc)

			resp, err := api.PayOrder(ctx, tc.args.req)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, &paymentv1.PayOrderResponse{}, resp)
				assert.NotEmpty(t, resp.GetTransactionUuid())
				assert.Equal(t, resp.GetTransactionUuid(), tc.expected.transactionUUID)
			}
		})
	}
}
