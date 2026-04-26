package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Create(ctx context.Context, order model.Order) error
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethodString) (uuid.UUID, error)
}
