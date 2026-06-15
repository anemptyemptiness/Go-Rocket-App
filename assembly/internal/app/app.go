package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/logger"
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
	}

	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
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
