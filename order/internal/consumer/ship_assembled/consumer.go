package ship_assembled

import (
	"context"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
)

type service struct {
	shipAssembledConsumer Consumer
	orderService          orderapi.OrderService
	orderRepository       ordersvc.OrderRepository
	inventoryClient       ordersvc.InventoryClient
}

func New(
	shipAssembledConsumer Consumer,
	orderService orderapi.OrderService,
	orderRepository ordersvc.OrderRepository,
	inventoryClient ordersvc.InventoryClient,
) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		orderService:          orderService,
		orderRepository:       orderRepository,
		inventoryClient:       inventoryClient,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	return s.shipAssembledConsumer.Consume(ctx, s.ShipAssembled)
}
