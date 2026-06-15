package assembly

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
)

type ShipAssembledProducerService interface {
	Produce(ctx context.Context, event model.ShipAssembledEvent) error
}

type Sleeper interface {
	Sleep() int64
}
