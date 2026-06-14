package config

type assemblerConfig struct {
	AssembleLimitTimeSec int64 `yaml:"limit" env-default:"5"`
}
