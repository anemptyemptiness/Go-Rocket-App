package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	authclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/iam/v1"
	inventoryclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1"
	paymentclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/middleware"
	orderrepository "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/order"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func NewHTTPHandlerWithProducer(
	orderPool *pgxpool.Pool,
	txManager orderservice.TxManager,
	inventoryClientGRPC inventoryv1.InventoryServiceClient,
	paymentClientGRPC paymentv1.PaymentServiceClient,
	authClientGRPC authv1.AuthServiceClient,
	orderPaidProducer orderservice.OrderPaidProducerService,
) (http.Handler, error) {
	paymentClient := paymentclientv1.New(paymentClientGRPC)
	inventoryClient := inventoryclientv1.New(inventoryClientGRPC)
	iamClient := authclientv1.New(authClientGRPC)

	orderRepo := orderrepository.New(orderPool)
	orderSvc := orderservice.New(
		orderRepo,
		paymentClient,
		inventoryClient,
		txManager,
		orderPaidProducer,
	)
	api := orderapi.NewAPI(orderSvc)

	handler, err := orderv1.NewServer(api)
	if err != nil {
		return nil, err
	}

	return middleware.AuthMiddleware(iamClient)(handler), nil
}
