package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethodString) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}
