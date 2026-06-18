package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderpaidproducersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/producer/order_producer"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/producer/order_producer/mocks"
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
		unexpectedErr = assert.AnError
	)

	validEvent := model.NewOrderPaidEvent(eventUUID, orderUUID, userUUID)
	validProto := eventsv1.OrderPaid{
		OrderUuid: validEvent.OrderUUID(),
		EventUuid: validEvent.EventUUID(),
		UserUuid:  validEvent.UserUUID(),
	}
	payload, err := proto.Marshal(&validProto)
	require.NoError(t, err)

	tests := []struct {
		name      string
		wantErr   bool
		args      model.OrderPaidEvent
		setupMock func(producer *mocks.Producer)
	}{
		{
			name:    "успешная обработки оплаченной детали",
			wantErr: false,
			args:    validEvent,
			setupMock: func(producer *mocks.Producer) {
				producer.EXPECT().
					Send(ctx, &kafka.Message{
						Key:   []byte(validEvent.EventUUID()),
						Value: payload,
					}).Return(nil)
			},
		},
		{
			name:    "ошибка отправки сообщения брокера сообщений",
			wantErr: true,
			args:    validEvent,
			setupMock: func(producer *mocks.Producer) {
				producer.EXPECT().
					Send(ctx, &kafka.Message{
						Key:   []byte(validEvent.EventUUID()),
						Value: payload,
					}).Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			producer := mocks.NewProducer(t)

			tc.setupMock(producer)

			producerSvc := orderpaidproducersvc.New(producer)

			producerErr := producerSvc.Produce(ctx, tc.args)
			if tc.wantErr {
				require.Error(t, producerErr)
			} else {
				require.NoError(t, producerErr)
			}
		})
	}
}
