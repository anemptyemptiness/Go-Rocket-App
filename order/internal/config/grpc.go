package config

import "net"

type grpcConfig struct {
	PaymentClient   paymentClientConfig   `yaml:"payment_client"`
	InventoryClient inventoryClientConfig `yaml:"inventory_client"`
	IAMClient       iamClientConfig       `yaml:"iam_client"`
}

type paymentClientConfig struct {
	Host string `yaml:"host" env:"PAYMENT_CLIENT_GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"PAYMENT_CLIENT_GRPC_PORT" env-default:"50052"`
}

func (p *paymentClientConfig) Address() string {
	return net.JoinHostPort(p.Host, p.Port)
}

type inventoryClientConfig struct {
	Host string `yaml:"host" env:"INVENTORY_CLIENT_GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"INVENTORY_CLIENT_GRPC_PORT" env-default:"50051"`
}

func (i *inventoryClientConfig) Address() string {
	return net.JoinHostPort(i.Host, i.Port)
}

type iamClientConfig struct {
	Host string `yaml:"host" env:"IAM_CLIENT_GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"IAM_CLIENT_GRPC_PORT" env-default:"50053"`
}

func (iam *iamClientConfig) Address() string {
	return net.JoinHostPort(iam.Host, iam.Port)
}
