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
		slog.ErrorContext(ctx, "консюмер order paid",
			slog.Int64("partition", int64(msg.Partition)),
			slog.Int64("offset", msg.Offset),
			slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "начинаем сборку корабля",
		slog.String("order_uuid", event.OrderUUID()),
		slog.Int64("offset", msg.Offset),
		slog.Int("partition", int(msg.Partition)),
	)

	err = s.assembleService.ShipAssemble(ctx, model.NewShipAssembledEvent(
		event.EventUUID(),
		event.OrderUUID(),
		event.UserUUID(),
	),
	)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "заказ собран успешно", "event_uuid", event.EventUUID())

	return nil
}
