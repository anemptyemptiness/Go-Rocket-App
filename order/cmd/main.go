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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	apiv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	inventoryclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1"
	paymentclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/order"
	orderService "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
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

	paymentClient := paymentclientv1.NewClient(paymentConn)
	inventoryClient := inventoryclientv1.NewClient(inventoryConn)
	orderRepo := order.NewRepository()
	orderSvc := orderService.NewService(orderRepo, paymentClient, inventoryClient)
	api := apiv1.NewAPI(orderSvc)

	// Создать OpenAPI сервер
	orderServer, err := orderv1.NewServer(api)
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
