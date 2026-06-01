package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	inventoryapi "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/api/inventory/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/config"
	inventoryrepo "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/part"
	applicationsvcpart "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/application/part"
	domainsvcchecker "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/domain/compatibility_checker"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	pgPool               *pgxpool.Pool
	inventoryAPI         inventoryv1.InventoryServiceServer
	inventoryService     inventoryapi.InventoryService
	compatibilityChecker applicationsvcpart.CompatibilityChecker
	inventoryRepo        applicationsvcpart.Repository
	txManager            applicationsvcpart.TxManager
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("создание пула соединений", "error", err)
			os.Exit(1)
		}

		closer.Add("pgPool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("проверка соединения с БД", "error", err)
			os.Exit(1)
		}

		d.pgPool = pool

		slog.Info("подключение к PostgreSQL установлено")
	}

	return d.pgPool
}

func (d *diContainer) InventoryAPI(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryAPI == nil {
		d.inventoryAPI = inventoryapi.New(d.InventoryService(ctx))
	}

	return d.inventoryAPI
}

func (d *diContainer) InventoryService(ctx context.Context) inventoryapi.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = applicationsvcpart.New(d.InventoryRepo(ctx), d.CompatibilityChecker(ctx), d.TxManager(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) CompatibilityChecker(_ context.Context) applicationsvcpart.CompatibilityChecker {
	if d.compatibilityChecker == nil {
		d.compatibilityChecker = domainsvcchecker.New()
	}

	return d.compatibilityChecker
}

func (d *diContainer) InventoryRepo(ctx context.Context) applicationsvcpart.Repository {
	if d.inventoryRepo == nil {
		d.inventoryRepo = inventoryrepo.New(d.PGPool(ctx))
	}

	return d.inventoryRepo
}

func (d *diContainer) TxManager(ctx context.Context) applicationsvcpart.TxManager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("создание transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = txManager

		slog.Info("transaction manager успешно создан")
	}

	return d.txManager
}
