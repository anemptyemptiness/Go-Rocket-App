package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	inventoryclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1"
	paymentclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1"
	orderrepository "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/order"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(
	orderPool *pgxpool.Pool,
	txManager orderrepository.TxManager,
	inventoryClientGRPC inventoryv1.InventoryServiceClient,
	paymentClientGRPC paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	paymentClient := paymentclientv1.New(paymentClientGRPC)
	inventoryClient := inventoryclientv1.New(inventoryClientGRPC)

	orderRepo := orderrepository.New(orderPool, txManager)
	orderSvc := orderservice.New(
		orderRepo,
		paymentClient,
		inventoryClient,
	)
	api := orderapi.NewAPI(orderSvc)

	handler, err := orderv1.NewServer(api)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
