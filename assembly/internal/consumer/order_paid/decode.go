package order_paid

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(_ context.Context, msg []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(msg, &pb); err != nil {
		return nil, err
	}

	return model.NewOrderPaidEvent(pb.GetEventUuid(), pb.GetOrderUuid(), pb.GetUserUuid()), nil
}
