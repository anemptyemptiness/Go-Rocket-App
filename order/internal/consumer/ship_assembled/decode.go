package ship_assembled

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

func decodeShipAssembled(_ context.Context, msg []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(msg, &pb); err != nil {
		return nil, err
	}

	return model.NewShipAssembledEvent(
		pb.GetEventUuid(),
		pb.GetOrderUuid(),
		pb.GetUserUuid(),
	), nil
}
