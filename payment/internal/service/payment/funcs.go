package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) PayOrder(ctx context.Context, req model.PayRequest) (string, error) {
	if req.OrderUUID == "" {
		slog.ErrorContext(ctx, "оплата заказа", slog.String("error", errs.ErrOrderUUIDIsEmpty.Error()))
		return "", pkgerr.InvalidArgument(errs.ErrOrderUUIDIsEmpty)
	}

	orderUuid, err := uuid.Parse(req.OrderUUID)
	if err != nil || orderUuid == uuid.Nil {
		slog.ErrorContext(ctx, "оплата заказа", slog.String("error", errs.ErrIncorrectOrderUUID.Error()))
		return "", pkgerr.InvalidArgument(errs.ErrIncorrectOrderUUID)
	}

	if req.PaymentMethod == model.PaymentMethodUnspecified {
		slog.ErrorContext(ctx, "оплата заказа", slog.String("error", errs.ErrPaymentMethodUnspecified.Error()))
		return "", pkgerr.InvalidArgument(errs.ErrPaymentMethodUnspecified)
	}

	transactionUUID := uuid.New()

	slog.InfoContext(ctx, "оплата прошла успешно",
		"order_uuid", req.OrderUUID,
		"transaction_uuid", transactionUUID.String(),
	)

	return transactionUUID.String(), nil
}
