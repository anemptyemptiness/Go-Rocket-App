package config

import "net"

type iamClientConfig struct {
	Host string `yaml:"host" env:"IAM_CLIENT_GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"IAM_CLIENT_GRPC_PORT" env-default:"50053"`
}

func (iam *iamClientConfig) Address() string {
	return net.JoinHostPort(iam.Host, iam.Port)
}
