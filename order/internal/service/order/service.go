package order

type service struct {
	orderRepository Repository
	paymentClient   PaymentClient
	inventoryClient InventoryClient
}

func NewService(
	orderRepository Repository,
	paymentClient PaymentClient,
	inventoryClient InventoryClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
