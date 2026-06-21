package assembly_consumer

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}
