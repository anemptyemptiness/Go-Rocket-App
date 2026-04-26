package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
}
