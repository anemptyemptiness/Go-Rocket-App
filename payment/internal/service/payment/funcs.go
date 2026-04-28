package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, orderUUIDStr string, paymentMethod model.PaymentMethod) (string, error) {
	if orderUUIDStr == "" {
		return "", errs.ErrOrderUUIDIsEmpty
	}

	orderUUID, err := uuid.Parse(orderUUIDStr)
	if err != nil {
		return "", errs.ErrIncorrectOrderUUID
	}

	if paymentMethod == model.PaymentMethodUnspecified {
		return "", errs.ErrPaymentMethodUnspecified
	}

	transactionUUID := uuid.New()

	slog.Info("оплата прошла успешно",
		"order_uuid", orderUUID,
		"transaction_uuid", transactionUUID.String(),
	)

	return transactionUUID.String(), nil
}
