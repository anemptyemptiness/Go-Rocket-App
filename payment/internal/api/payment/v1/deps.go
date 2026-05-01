package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req model.PayRequest) (string, error)
}
