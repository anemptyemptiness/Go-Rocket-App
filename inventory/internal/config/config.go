package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var appConfig *config

type config struct {
	GRPC   grpcConfig   `yaml:"grpc"`
	Logger loggerConfig `yaml:"logger"`
	PG     pgConfig     `yaml:"pg"`
}

const defaultConfigPath = "config.local.yaml"

// ResolveConfigPath определяет путь к конфиг-файлу по цепочке приоритетов:
// флаг -config > env CONFIG_PATH > "config.local.yaml".
func ResolveConfigPath() string {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "путь к YANL-конфигу")
	flag.Parse()

	if cfgPath != "" {
		return cfgPath
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return defaultConfigPath
}

func Load(path string) error {
	var cfg config

	if path != "" {
		// ReadConfig читает YAML, затем перетирает значения из env.
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return fmt.Errorf("загрузить конфиг из %s: %w", path, err)
		}

		appConfig = &cfg

		return nil
	}

	// Если путь не указан — только env-переменные.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return fmt.Errorf("загрузить конфиг из env: %w", err)
	}

	appConfig = &cfg

	return nil
}

func AppConfig() *config {
	return appConfig
}
