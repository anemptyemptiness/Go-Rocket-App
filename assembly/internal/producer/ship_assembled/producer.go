package ship_assembled

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

type service struct {
	shipAssembledProducer Producer
}

func New(shipAssembledProducer Producer) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) Produce(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID(),
		OrderUuid:    event.OrderUUID(),
		UserUuid:     event.UserUUID(),
		BuildTimeSec: event.BuildTimeSec(),
		AssembledAt:  timestamppb.New(event.AssembledAt()),
	}

	payload, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	err = s.shipAssembledProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.EventUUID()),
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}
