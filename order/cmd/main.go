package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/anemptyemptiness/Go-Rocket-App/order/pkg/app"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "0.0.0.0:50051"
	paymentServiceAddress   = "0.0.0.0:50052"

	httpHost = "localhost"
	httpPort = "8080"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20 // 1 Mb

	shutdownTimeout = 10 * time.Second

	keepAliveTime                = 60 * time.Second
	keepAliveTimeout             = 3 * time.Second
	keepAlivePermitWithoutStream = true
)

func main() {
	ctx := context.Background()

	err := godotenv.Load("../order.env")
	if err != nil {
		slog.Error("не удалось загрузить окружение", "error", err)
		return
	}

	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: keepAlivePermitWithoutStream,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
	}
	defer inventoryConn.Close()
	inventoryClientGRPC := inventoryv1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
	}
	defer paymentConn.Close()
	paymentClientGRPC := paymentv1.NewPaymentServiceClient(paymentConn)

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

	// Создать OpenAPI сервер
	orderServer, err := app.NewHTTPHandler(pool, txManager, inventoryClientGRPC, paymentClientGRPC)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
	}

	httpServer := &http.Server{
		Addr:              net.JoinHostPort(httpHost, httpPort),
		Handler:           orderServer,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	slog.Info("запуск OrderService", "port", 8080)

	go func() {
		err = httpServer.ListenAndServe()
		if err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown http-сервера
	slog.Info("🛑 завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := httpServer.Shutdown(ctx); shutdownErr != nil {
		slog.Error("❌ ошибка при остановке сервера", "error", shutdownErr)
	}

	slog.Info("✅ сервер остановлен")
}
