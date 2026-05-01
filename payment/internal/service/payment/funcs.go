package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, req model.PayRequest) (string, error) {
	if req.OrderUUID == "" {
		return "", errs.ErrOrderUUIDIsEmpty
	}

	orderUuid, err := uuid.Parse(req.OrderUUID)
	if err != nil || orderUuid == uuid.Nil {
		return "", errs.ErrIncorrectOrderUUID
	}

	if req.PaymentMethod == model.PaymentMethodUnspecified {
		return "", errs.ErrPaymentMethodUnspecified
	}

	transactionUUID := uuid.New()

	slog.Info("оплата прошла успешно",
		"order_uuid", req.OrderUUID,
		"transaction_uuid", transactionUUID.String(),
	)

	return transactionUUID.String(), nil
}
