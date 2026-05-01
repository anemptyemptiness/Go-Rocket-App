package app

import (
	"google.golang.org/grpc"

	paymentapi "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/interceptor"
	paymentsvc "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/service/payment"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	svc := paymentsvc.New()
	api := paymentapi.New(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorInterceptor(),
		),
	}

	return opts
}
