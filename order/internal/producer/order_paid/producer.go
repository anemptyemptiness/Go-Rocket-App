package order_paid

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer Producer
}

func New(orderPaidProducer Producer) ordersvc.OrderPaidProducerService {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) Produce(ctx context.Context, event model.OrderPaidEvent) error {
	msg := eventsv1.OrderPaid{
		OrderUuid: event.OrderUUID(),
		EventUuid: event.EventUUID(),
		UserUuid:  event.UserUUID(),
	}

	payload, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	err = s.orderPaidProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.EventUUID()),
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}
