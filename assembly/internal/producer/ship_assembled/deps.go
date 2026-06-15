package ship_assembled

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
)

type Producer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
