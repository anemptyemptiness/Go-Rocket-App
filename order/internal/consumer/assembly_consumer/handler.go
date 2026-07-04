package assembly_consumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
)

func (s *service) ShipAssembled(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(ctx, msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "консюмер ship assembled", slog.String("error", err.Error()))
		return err
	}

	order, err := s.orderRepository.Get(ctx, event.OrderUUID())
	if err != nil {
		slog.ErrorContext(ctx, "консюмер ship assembled", slog.String("error", err.Error()))
		return err
	}

	if order.Status == model.OrderStatusAssembled {
		slog.InfoContext(ctx, "консюмер ship assembled: заказ уже в статусе "+string(model.OrderStatusAssembled))
		return nil
	}

	inventoryClientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	partUUIDs := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		partUUIDs = append(partUUIDs, item.PartUuid)
	}

	err = s.inventoryClient.CommitParts(inventoryClientCtx, partUUIDs)
	if err != nil {
		slog.ErrorContext(ctx, "консюмер ship assembled", slog.String("error", err.Error()))
		return err
	}

	order.SetStatus(model.OrderStatusAssembled)

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		slog.ErrorContext(ctx, "консюмер ship assembled", slog.String("error", err.Error()))
		return err
	}

	return nil
}
