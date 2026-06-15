package assembly

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
)

func (s *service) ShipAssemble(ctx context.Context, event model.ShipAssembledEvent) error {
	buildTimeSec := s.sleeper.Sleep()

	event.SetBuildTimeSec(buildTimeSec)
	event.MarkAssembledAt()

	err := s.shipAssembledProducerSvc.Produce(ctx, event)
	if err != nil {
		slog.Error("не удалось опубликовать событие ShipAssembled в Kafka", "error", err)
		return fmt.Errorf("опубликовать событие ShipAssembled: %w", err)
	}

	return nil
}
