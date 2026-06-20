package assembly_consumer

import (
	"context"

	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
)

type service struct {
	shipAssembledConsumer Consumer
	orderRepository       ordersvc.OrderRepository
	inventoryClient       ordersvc.InventoryClient
}

func New(
	shipAssembledConsumer Consumer,
	orderRepository ordersvc.OrderRepository,
	inventoryClient ordersvc.InventoryClient,
) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		orderRepository:       orderRepository,
		inventoryClient:       inventoryClient,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	return s.shipAssembledConsumer.Consume(ctx, s.ShipAssembled)
}
