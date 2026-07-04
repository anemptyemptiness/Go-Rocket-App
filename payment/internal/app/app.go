package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/grpc/health"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/logger"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/metrics"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/tracing"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

const (
	maxConnectionIdle     = 15 * time.Minute
	maxConnectionAge      = 30 * time.Minute
	maxConnectionAgeGrace = 5 * time.Second

	keepAliveTime                = 5 * time.Minute
	keepAliveTimeout             = 10 * time.Second
	keepAliveMinTime             = 50 * time.Second
	keepAlivePermitWithoutStream = true

	shutdownTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	listener    net.Listener
	grpcServer  *grpc.Server
}

func New(ctx context.Context) *App {
	app := &App{}

	app.initDeps(ctx)

	return app
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) {
	deps := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range deps {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(logger.Config{
		Level:             config.AppConfig().Logger.Level,
		ServiceName:       config.AppConfig().OTel.ServiceName,
		Environment:       config.AppConfig().OTel.Environment,
		EnableOTLP:        config.AppConfig().OTel.EnableOTLP,
		CollectorEndpoint: config.AppConfig().OTel.Endpoint,
	})

	closer.Add("logger", func(context.Context) error {
		loggerCloseErr := logger.Close()
		if loggerCloseErr != nil {
			return loggerCloseErr
		}
		return nil
	})
}

func (a *App) initMetrics(_ context.Context) {
	metrics.Init(config.AppConfig().OTel.ServiceName)

	closer.Add("metrics", func(context.Context) error {
		metricsCloseErr := metrics.Close()
		if metricsCloseErr != nil {
			return metricsCloseErr
		}
		return nil
	})
}

func (a *App) initTracing(ctx context.Context) {
	shutdown, err := tracing.InitTracer(ctx, tracing.Config{
		CollectorEndpoint: config.AppConfig().OTel.Endpoint,
		ServiceName:       config.AppConfig().OTel.ServiceName,
		Environment:       config.AppConfig().OTel.Environment,
		ServiceVersion:    config.AppConfig().OTel.ServiceVersion,
		SamplingRatio:     config.AppConfig().OTel.SamplingRatio,
	})
	if err != nil {
		slog.Error("ошибка при инициализации трейсера", "error", err)
		os.Exit(1)
	}

	closer.Add("tracing", func(closerCtx context.Context) error {
		return shutdown(closerCtx)
	})
}

func (a *App) initListener(ctx context.Context) {
	lc := net.ListenConfig{}

	listener, err := lc.Listen(ctx, "tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}

func (a *App) initGRPCServer(_ context.Context) {
	a.grpcServer = grpc.NewServer(
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
		grpc.UnaryInterceptor(pkgerr.UnaryErrorInterceptor(slog.Default())),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	api := a.diContainer.PaymentAPI()

	closer.Add("gRPC Server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	paymentv1.RegisterPaymentServiceServer(a.grpcServer, api)
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() { //nolint:gosec // G118: ctx уже отменён, context.Background нужен для graceful shutdown
		<-ctx.Done()

		cancel()

		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			slog.Error("ошибка при завершении работы", "error", closeErr)
		}
	}()
}

func (a *App) runGRPCServer() error {
	slog.Info("🚀 gRPC-сервер запущен", "address", config.AppConfig().GRPC.Address())

	return a.grpcServer.Serve(a.listener)
}
