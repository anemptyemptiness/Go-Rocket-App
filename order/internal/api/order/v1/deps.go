package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, orderUUID string) error
}
