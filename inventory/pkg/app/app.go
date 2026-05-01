package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/interceptor"
	partrepo "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/part"
	partsvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/part"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool, manager partsvc.TxManager) {
	repo := partrepo.New(pool)
	svc := partsvc.New(repo, manager)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorInterceptor(),
		),
	}

	return opts
}
