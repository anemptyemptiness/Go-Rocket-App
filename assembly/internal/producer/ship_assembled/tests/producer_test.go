package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	shipassembledsvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/producer/ship_assembled"
	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/producer/ship_assembled/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

func Test_Produce(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		eventUUID     = gofakeit.UUID()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		buildTimeSec  = int64(1)
		unexpectedErr = assert.AnError
	)

	validEvent := model.NewShipAssembledEvent(eventUUID, orderUUID, userUUID)
	validEvent.SetBuildTimeSec(buildTimeSec)
	validEvent.MarkAssembledAt()

	tests := []struct {
		name      string
		wantErr   bool
		args      model.ShipAssembledEvent
		setupMock func(producer *mocks.Producer)
	}{
		{
			name:    "успешная обработка сообщения",
			wantErr: false,
			args:    validEvent,
			setupMock: func(producer *mocks.Producer) {
				producer.EXPECT().
					Send(ctx, mock.MatchedBy(func(msg *kafka.Message) bool {
						if string(msg.Key) != eventUUID {
							return false
						}

						var event eventsv1.ShipAssembled
						if err := proto.Unmarshal(msg.Value, &event); err != nil {
							return false
						}

						return event.GetEventUuid() == validEvent.EventUUID() &&
							event.GetOrderUuid() == validEvent.OrderUUID() &&
							event.GetUserUuid() == validEvent.UserUUID() &&
							event.GetBuildTimeSec() == validEvent.BuildTimeSec() &&
							event.GetAssembledAt() != nil &&
							event.GetAssembledAt().AsTime().Equal(validEvent.AssembledAt())
					})).
					Return(nil)
			},
		},
		{
			name:    "ошибка отправки сообщения в продюсере",
			wantErr: true,
			args:    validEvent,
			setupMock: func(producer *mocks.Producer) {
				producer.EXPECT().
					Send(ctx, mock.MatchedBy(func(msg *kafka.Message) bool {
						if string(msg.Key) != eventUUID {
							return false
						}

						var event eventsv1.ShipAssembled
						if err := proto.Unmarshal(msg.Value, &event); err != nil {
							return false
						}

						return event.GetEventUuid() == validEvent.EventUUID() &&
							event.GetOrderUuid() == validEvent.OrderUUID() &&
							event.GetUserUuid() == validEvent.UserUUID() &&
							event.GetBuildTimeSec() == validEvent.BuildTimeSec() &&
							event.GetAssembledAt() != nil &&
							event.GetAssembledAt().AsTime().Equal(validEvent.AssembledAt())
					})).
					Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			producer := mocks.NewProducer(t)

			tc.setupMock(producer)

			svc := shipassembledsvc.New(producer)

			produceErr := svc.Produce(ctx, tc.args)
			if tc.wantErr {
				require.Error(t, produceErr)
			} else {
				require.NoError(t, produceErr)
			}
		})
	}
}
