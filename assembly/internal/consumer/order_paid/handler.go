package order_paid

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
)

func (s *service) OrderPaid(ctx context.Context, msg kafka.Message) error {
	event, err := decodeOrderPaid(ctx, msg.Value)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "заказ оплачен", "event_uuid", event.EventUUID())

	err = s.assembleService.ShipAssemble(ctx, model.NewShipAssembledEvent(
		event.EventUUID(),
		event.OrderUUID(),
		event.UserUUID(),
	),
	)
	if err != nil {
		return err
	}

	return nil
}
