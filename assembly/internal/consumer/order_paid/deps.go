package order_paid

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
)

type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

type AssembleService interface {
	ShipAssemble(ctx context.Context, event model.ShipAssembledEvent) error
}
