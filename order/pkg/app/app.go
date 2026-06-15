package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	inventoryclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1"
	paymentclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1"
	orderpaidproducer "github.com/anemptyemptiness/Go-Rocket-App/order/internal/producer/order_paid"
	orderrepository "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/order"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/producer"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(
	orderPool *pgxpool.Pool,
	txManager orderservice.TxManager,
	inventoryClientGRPC inventoryv1.InventoryServiceClient,
	paymentClientGRPC paymentv1.PaymentServiceClient,
	orderPaidProducer *producer.Producer,
) (http.Handler, error) {
	paymentClient := paymentclientv1.New(paymentClientGRPC)
	inventoryClient := inventoryclientv1.New(inventoryClientGRPC)
	p := orderpaidproducer.New(orderPaidProducer)

	orderRepo := orderrepository.New(orderPool)
	orderSvc := orderservice.New(
		orderRepo,
		paymentClient,
		inventoryClient,
		txManager,
		p,
	)
	api := orderapi.NewAPI(orderSvc)

	handler, err := orderv1.NewServer(api)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
