package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/interceptor"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/pkg/app"
)

const (
	grpcAddress = "0.0.0.0:50051"

	maxConnectionIdle     = 15 * time.Minute
	maxConnectionAge      = 30 * time.Minute
	maxConnectionAgeGrace = 5 * time.Second

	keepAliveTime                = 5 * time.Minute
	keepAliveTimeout             = 10 * time.Second
	keepAliveMinTime             = 50 * time.Second
	keepAlivePermitWithoutStream = true
)

func main() {
	ctx := context.Background()

	err := godotenv.Load("../inventory.env")
	if err != nil {
		slog.Error("не удалось загрузить окружение", "error", err)
		return
	}

	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     maxConnectionIdle,
			MaxConnectionAge:      maxConnectionAge,
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
			Time:                  keepAliveTime,
			Timeout:               keepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepAliveMinTime,
			PermitWithoutStream: keepAlivePermitWithoutStream,
		}),
		grpc.ChainUnaryInterceptor(
			interceptor.ErrorInterceptor(),
		),
	)

	// DSN берём из order.env / inventory.env (пока хардкодим в main.go, конфиги — неделя 4)
	pool, err := pgxpool.New(ctx, os.Getenv("DB_URI"))
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
	}
	defer pool.Close()

	// Проверяем соединение
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("проверка соединения с БД", "error", err)
	}

	slog.Info("подключение к PostgreSQL установлено")

	// Создаём Transaction Manager для pgx
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("создание transaction manager", "error", err)
	}

	app.RegisterServices(grpcServer, pool, txManager)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 остановка gRPC сервера")

	grpcServer.GracefulStop()

	slog.Info("✅ сервер остановлен")
}
