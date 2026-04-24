package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.CreateOrderResponse, error)
	PayOrder(ctx context.Context, req model.PayOrderRequest, orderUUID uuid.UUID) (model.PayOrderResponse, error)
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
}
