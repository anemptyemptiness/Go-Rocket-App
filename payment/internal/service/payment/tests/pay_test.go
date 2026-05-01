package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/service/payment"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req model.PayRequest
	}

	type expected struct {
		err error
	}

	var (
		ctx            = context.Background()
		orderUUID      = gofakeit.UUID()
		emptyOrderUUID = ""
	)

	tests := []struct {
		name       string
		args       args
		expected   expected
		setupMocks func()
	}{
		{
			name: "успешная оплата картой",
			args: args{
				req: model.PayRequest{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			expected: expected{
				err: nil,
			},
			setupMocks: func() {},
		},
		{
			name: "успешная оплата СБП",
			args: args{
				req: model.PayRequest{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
			expected: expected{
				err: nil,
			},
			setupMocks: func() {},
		},
		{
			name: "успешная оплата кредитной картой",
			args: args{
				req: model.PayRequest{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodCreditCard,
				},
			},
			expected: expected{
				err: nil,
			},
			setupMocks: func() {},
		},
		{
			name: "успешная оплата деньгами инвестора",
			args: args{
				req: model.PayRequest{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodInvestorMoney,
				},
			},
			expected: expected{
				err: nil,
			},
			setupMocks: func() {},
		},
		{
			name: "ошибка: orderUUID пустой",
			args: args{
				req: model.PayRequest{
					OrderUUID:     emptyOrderUUID,
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
			expected: expected{
				err: errs.ErrOrderUUIDIsEmpty,
			},
			setupMocks: func() {},
		},
		{
			name: "ошибка: метод оплаты неопределённый",
			args: args{
				req: model.PayRequest{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodUnspecified,
				},
			},
			expected: expected{
				err: errs.ErrPaymentMethodUnspecified,
			},
			setupMocks: func() {},
		},
		{
			name: "ошибка: идентификатор заказа невалидный",
			args: args{
				req: model.PayRequest{
					OrderUUID:     "invalid!!!",
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
			expected: expected{
				err: errs.ErrIncorrectOrderUUID,
			},
			setupMocks: func() {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := payment.New()

			transactionStr, err := svc.PayOrder(ctx, tc.args.req)
			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, transactionStr)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, transactionStr)

				transactionUUID, parseErr := uuid.Parse(transactionStr)
				require.NoError(t, parseErr)
				assert.NotEqual(t, transactionUUID, uuid.Nil)
			}
		})
	}
}
