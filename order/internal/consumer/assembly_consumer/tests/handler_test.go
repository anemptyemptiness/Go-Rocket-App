package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	apimocks "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1/mocks"
	shipassembledconsumer "github.com/anemptyemptiness/Go-Rocket-App/order/internal/consumer/assembly_consumer"
	consumermocks "github.com/anemptyemptiness/Go-Rocket-App/order/internal/consumer/assembly_consumer/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	svcmocks "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka"
	eventsv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/events/v1"
)

func Test_ShipAssembled(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = gofakeit.UUID()
		eventUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()
		buildTimeSec    = int64(1)

		part1UUID = gofakeit.UUID()
		part2UUID = gofakeit.UUID()

		now           = time.Now()
		unexpectedErr = assert.AnError
	)

	orderNotAssembled := model.Order{
		UUID:     orderUUID,
		UserUUID: userUUID,
		Items: []model.OrderItem{
			{
				PartUuid: part1UUID,
			},
			{
				PartUuid: part2UUID,
			},
		},
		TotalPrice:      5000,
		TransactionUUID: new(transactionUUID),
		PaymentMethod:   new(model.PaymentMethodCard),
		Status:          model.OrderStatusPaid,
		CreatedAt:       now,
	}

	orderAlreadyAssembled := model.Order{
		UUID:     orderUUID,
		UserUUID: userUUID,
		Items: []model.OrderItem{
			{
				PartUuid: part1UUID,
			},
			{
				PartUuid: part2UUID,
			},
		},
		TotalPrice:      5000,
		TransactionUUID: new(transactionUUID),
		PaymentMethod:   new(model.PaymentMethodCard),
		Status:          model.OrderStatusAssembled,
		CreatedAt:       now,
	}

	validProto := eventsv1.ShipAssembled{
		EventUuid:    eventUUID,
		OrderUuid:    orderUUID,
		UserUuid:     userUUID,
		BuildTimeSec: buildTimeSec,
		AssembledAt:  timestamppb.New(now),
	}

	validProtoMarshalled, marshalErr := proto.Marshal(&validProto)
	require.NoError(t, marshalErr)

	tests := []struct {
		name      string
		wantErr   bool
		args      kafka.Message
		setupMock func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient)
	}{
		{
			name:    "успешная сборка детали",
			wantErr: false,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, mock.MatchedBy(func(eventOrderUUID string) bool {
						return orderUUID == eventOrderUUID
					})).Return(orderNotAssembled, nil)

				invClient.EXPECT().
					CommitParts(mock.Anything, []string{part1UUID, part2UUID}).
					Return(nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.Status == model.OrderStatusAssembled
					})).
					Return(nil)
			},
		},
		{
			name:    "ошибка репозитория (Get)",
			wantErr: true,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, mock.MatchedBy(func(eventOrderUUID string) bool {
						return orderUUID == eventOrderUUID
					})).
					Return(model.Order{}, unexpectedErr)
			},
		},
		{
			name:    "деталь уже собрана",
			wantErr: false,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, mock.MatchedBy(func(eventOrderUUID string) bool {
						return orderUUID == eventOrderUUID
					})).Return(orderAlreadyAssembled, nil)
			},
		},
		{
			name:    "ошибка inventory клиента (CommitParts)",
			wantErr: true,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, mock.MatchedBy(func(eventOrderUUID string) bool {
						return orderUUID == eventOrderUUID
					})).Return(orderNotAssembled, nil)

				invClient.EXPECT().
					CommitParts(mock.Anything, []string{part1UUID, part2UUID}).
					Return(unexpectedErr)
			},
		},
		{
			name:    "ошибка репозитория (Update)",
			wantErr: true,
			args: kafka.Message{
				Key:   []byte(eventUUID),
				Value: validProtoMarshalled,
			},
			setupMock: func(repo *svcmocks.OrderRepository, invClient *svcmocks.InventoryClient) {
				repo.EXPECT().
					Get(ctx, mock.MatchedBy(func(eventOrderUUID string) bool {
						return orderUUID == eventOrderUUID
					})).Return(orderNotAssembled, nil)

				invClient.EXPECT().
					CommitParts(mock.Anything, []string{part1UUID, part2UUID}).
					Return(nil)

				repo.EXPECT().
					Update(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.Status == model.OrderStatusAssembled
					})).Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			consumer := consumermocks.NewConsumer(t)
			repo := svcmocks.NewOrderRepository(t)
			orderSvc := apimocks.NewOrderService(t)
			invClient := svcmocks.NewInventoryClient(t)

			svc := shipassembledconsumer.New(consumer, orderSvc, repo, invClient)

			tc.setupMock(repo, invClient)

			shipAssembledErr := svc.ShipAssembled(ctx, tc.args)
			if tc.wantErr {
				require.Error(t, shipAssembledErr)
			} else {
				require.NoError(t, shipAssembledErr)
			}
		})
	}
}
