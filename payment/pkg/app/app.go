package app

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentapi "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1"
	paymentsvc "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/service/payment"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

// PaymentServer реализует gRPC сервис оплаты.
type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func RegisterServices(grpcServer *grpc.Server) {
	svc := paymentsvc.New()
	api := paymentapi.New(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0)

	return opts
}

// PayOrder обрабатывает оплату заказа.
func (s *PaymentServer) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "идентификатор заказа не может быть пустым")
	}

	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "метод оплаты неопределённый")
	}

	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "идентификатор заказа невалиден: %v", err)
	}

	transactionUUID := uuid.New()

	slog.Info("оплата прошла успешно",
		"order_uuid", orderUUID.String(),
		"transaction_uuid", transactionUUID.String(),
	)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
