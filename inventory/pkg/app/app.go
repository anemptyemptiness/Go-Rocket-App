package app

import (
	"log/slog"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	partrepo "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/part"
	applicationpartsvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part"
	domainpartsvc "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/domain/compatibility_checker"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	// Создаём Transaction Manager для pgx.
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("создание transaction manager", "error", err)
		return
	}

	repo := partrepo.New(pool)
	compChecker := domainpartsvc.New()
	svc := applicationpartsvc.New(repo, compChecker, txManager)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}

func Interceptors() []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			pkgerr.UnaryErrorInterceptor(slog.Default()),
		),
	}

	return opts
}
