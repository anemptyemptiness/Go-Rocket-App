package app

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/logger"
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
		errChan <- a.runHTTPServer()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		a.startGracefulShutdown()
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
	} else {
		slog.Info("http server успешно остановлен")
	}
}

func (a *App) initDeps(ctx context.Context) {
	deps := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initHTTPServer,
	}

	for _, f := range deps {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initHTTPServer(ctx context.Context) {
	closer.Add("HTTP Server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			slog.Error("ошибка при остановке HTTP-сервера", "error", err)
		} else {
			slog.Info("HTTP-сервер успешно остановлен")
		}
		return nil
	})

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           a.diContainer.OrderServer(ctx),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}
