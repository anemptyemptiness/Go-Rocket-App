package config

// internal/config/otel.go
// otelConfig — внутренняя структура пакета config (имя с маленькой буквы),
// наружу торчит только метод OTel() OTelConfig интерфейса.
type otelConfig struct {
	EnableOTLP     bool    `yaml:"enable_otlp"`
	Endpoint       string  `yaml:"endpoint"` // адрес OTEL Collector (OTLP gRPC)
	Environment    string  `yaml:"environment"`
	ServiceName    string  `yaml:"service_name"` // имя сервиса в трейсах / метриках / логах
	ServiceVersion string  `yaml:"service_version"`
	SamplingRatio  float64 `yaml:"sampling_ratio"`
}
