package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/app"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/config"
)

func main() {
	// .env опционален — ошибка загрузки допустима.
	_ = godotenv.Load("../iam.env") //nolint:gosec // G104: ошибка загрузки допустима.

	configPath := config.ResolveConfigPath()

	// YAML-конфиг + env-переменные поверх.
	err := config.Load(configPath)
	if err != nil {
		slog.Error("не удалось загрузить конфигурацию", "error", err)
		os.Exit(1)
	}

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
