package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type Repository interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, req model.ListPartsClientRequest) (model.ListPartsClientResponse, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, req model.PayOrderClientRequest) (model.PayOrderClientResponse, error)
}
