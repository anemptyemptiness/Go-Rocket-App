package main

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/pkg/app"
	errs "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

const (
	grpcAddress = "0.0.0.0:50052"

	maxConnectionIdle     = 15 * time.Minute
	maxConnectionAge      = 30 * time.Minute
	maxConnectionAgeGrace = 5 * time.Second

	keepAliveTime                = 5 * time.Minute
	keepAliveTimeout             = 10 * time.Second
	keepAliveMinTime             = 50 * time.Second
	keepAlivePermitWithoutStream = true
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return
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
		grpc.UnaryInterceptor(errs.UnaryErrorInterceptor(slog.Default())),
	)

	app.RegisterServices(grpcServer)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск PaymentService", "адрес", grpcAddress)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
			return
		}
	}()

	<-ctx.Done()
	slog.Info("🛑 остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("✅ сервер остановлен")
}
