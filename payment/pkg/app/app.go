package app

import (
	"log/slog"

	"google.golang.org/grpc"

	paymentapi "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1"
	paymentsvc "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/service/payment"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
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
			pkgerr.UnaryErrorInterceptor(slog.Default()),
		),
	}

	return opts
}
