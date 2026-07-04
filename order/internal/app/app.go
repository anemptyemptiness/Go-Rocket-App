package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/middleware"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/logger"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/metrics"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/tracing"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20 // 1 Mb

	keepAliveTime                = 60 * time.Second
	keepAliveTimeout             = 3 * time.Second
	keepAlivePermitWithoutStream = true

	shutdownTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	consumer    ShipAssembledConsumer
}

func New(ctx context.Context) *App {
	app := &App{}

	app.initDeps(ctx)

	return app
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error, 2)
	go func() {
		errChan <- a.runHTTPServer()
	}()

	go func() {
		errChan <- a.runConsumer(ctx)
	}()

	select {
	case err := <-errChan:
		a.startGracefulShutdown()
		return err
	case <-ctx.Done():
		a.startGracefulShutdown()
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	slog.Info("🚀 Kafka-консюмер запущен", "address", config.AppConfig().Kafka.Brokers)

	if consumerErr := a.consumer.RunConsumer(ctx); consumerErr != nil {
		return consumerErr
	}
	return nil
}

func (a *App) runHTTPServer() error {
	slog.Info("🚀 HTTP-сервер запущен", "address", config.AppConfig().HTTP.Address())

	if httpErr := a.httpServer.ListenAndServe(); httpErr != nil && !errors.Is(httpErr, http.ErrServerClosed) {
		return httpErr
	}
	return nil
}

func (a *App) startGracefulShutdown() {
	slog.Info("получен сигнал завершения, начинаем graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
		slog.Error("ошибка при завершении работы", "error", closeErr)
	}
}

func (a *App) initDeps(ctx context.Context) {
	deps := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initHTTPServer,
		a.initShipAssembledConsumer,
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

func (a *App) initHTTPServer(ctx context.Context) {
	// Оборачиваем хэндлер HTTP в auth-мидлвари.
	handler := middleware.AuthMiddleware(a.diContainer.IAMClient(ctx))(a.diContainer.OrderServer(ctx))
	// Оборачиваем хэндлер HTTP в tracing-мидлвари.
	handler = otelhttp.NewHandler(handler, config.AppConfig().OTel.ServiceName)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	closer.Add("HTTP Server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			slog.Error("ошибка при остановке HTTP-сервера", "error", err)
		} else {
			slog.Info("HTTP-сервер успешно остановлен")
		}
		return nil
	})
}

func (a *App) initShipAssembledConsumer(ctx context.Context) {
	a.consumer = a.diContainer.ShipAssembledConsumer(ctx)
}
