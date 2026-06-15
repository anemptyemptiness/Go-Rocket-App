package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	consumerorderpaid "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/consumer/order_paid"
	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/consumer/order_paid/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

func Test_OrderPaid(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		eventUUID     = gofakeit.UUID()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	validProto := eventsv1.OrderPaid{
		EventUuid: eventUUID,
		OrderUuid: orderUUID,
		UserUuid:  userUUID,
	}

	validProtoMarshalled, marshalErr := proto.Marshal(&validProto)
	require.NoError(t, marshalErr)

	tests := []struct {
		name      string
		wantErr   bool
		args      kafka.Message
		setupMock func(svc *mocks.AssembleService)
	}{
		{
			name:    "успешная обработка сообщения",
			wantErr: false,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(svc *mocks.AssembleService) {
				svc.EXPECT().
					ShipAssemble(ctx, mock.MatchedBy(func(event model.ShipAssembledEvent) bool {
						return event.EventUUID() != "" && event.EventUUID() == eventUUID &&
							event.OrderUUID() != "" && event.OrderUUID() == orderUUID &&
							event.UserUUID() != "" && event.UserUUID() == userUUID
					})).Return(nil)
			},
		},
		{
			name:    "ошибка: невалидное proto-сообщение",
			wantErr: true,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: []byte("invalid proto message"),
			},
			setupMock: func(svc *mocks.AssembleService) {},
		},
		{
			name:    "ошибка: сервис вернул ошибку",
			wantErr: true,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(svc *mocks.AssembleService) {
				svc.EXPECT().
					ShipAssemble(ctx, mock.MatchedBy(func(event model.ShipAssembledEvent) bool {
						return event.EventUUID() != "" && event.EventUUID() == eventUUID &&
							event.OrderUUID() != "" && event.OrderUUID() == orderUUID &&
							event.UserUUID() != "" && event.UserUUID() == userUUID
					})).Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			consumer := mocks.NewConsumer(t)
			svc := mocks.NewAssembleService(t)

			tc.setupMock(svc)

			wrappedConsumer := consumerorderpaid.New(consumer, svc)

			orderPaidErr := wrappedConsumer.OrderPaid(ctx, tc.args)
			if tc.wantErr {
				require.Error(t, orderPaidErr)
			} else {
				require.NoError(t, orderPaidErr)
			}
		})
	}
}
