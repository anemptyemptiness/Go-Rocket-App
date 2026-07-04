package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/logger"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/metrics"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/tracing"
)

const (
	shutdownTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) *App {
	app := &App{}

	app.initDeps(ctx)

	return app
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		err := a.runConsumer(ctx)
		if err != nil {
			slog.Error("не удалось запустить Kafka-потребитель", "error", err)
			errChan <- err
			return
		}

		errChan <- nil
	}()

	var runErr error
	select {
	case runErr = <-errChan:
	case <-ctx.Done():
	}

	gracefulErr := a.startGracefulShutdown(ctx, cancel)
	if gracefulErr != nil {
		if runErr == nil {
			runErr = gracefulErr
		}
	}

	return runErr
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
	}

	for _, f := range inits {
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

func (a *App) runConsumer(ctx context.Context) error {
	slog.Info("🚀 Kafka-потребитель OrderPaid запущен")

	return a.diContainer.OrderPaidConsumer().RunConsumer(ctx)
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) error {
	<-ctx.Done()

	cancel()

	slog.Info("получен сигнал завершения, начинаем graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
		slog.Error("ошибка при завершении работы", "error", closeErr)
		return closeErr
	}

	return nil
}
