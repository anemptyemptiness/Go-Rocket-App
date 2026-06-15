package order

type service struct {
	orderRepository   OrderRepository
	paymentClient     PaymentClient
	inventoryClient   InventoryClient
	txManager         TxManager
	orderPaidProducer OrderPaidProducerService
}

func New(
	orderRepository OrderRepository,
	paymentClient PaymentClient,
	inventoryClient InventoryClient,
	txManager TxManager,
	orderPaidProducer OrderPaidProducerService,
) *service {
	return &service{
		orderRepository:   orderRepository,
		paymentClient:     paymentClient,
		inventoryClient:   inventoryClient,
		txManager:         txManager,
		orderPaidProducer: orderPaidProducer,
	}
}
