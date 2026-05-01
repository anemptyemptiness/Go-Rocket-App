package order

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (string, error)
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
