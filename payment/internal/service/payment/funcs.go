package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, req model.PayOrderRequest) (model.PayOrderResponse, error) {
	transactionUUID := uuid.New()

	slog.Info("оплата прошла успешно",
		"order_uuid", req.OrderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
