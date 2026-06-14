package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/model"
	assemblysvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/service/assembly"
	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/service/assembly/mocks"
)

func Test_ShipAssemble(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		eventUUID     = gofakeit.UUID()
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		buildTimeSec  = int64(1)
		unexpectedErr = assert.AnError
	)

	event := model.NewShipAssembledEvent(eventUUID, orderUUID, userUUID)

	tests := []struct {
		name      string
		wantErr   bool
		args      model.ShipAssembledEvent
		setupMock func(sleeper *mocks.Sleeper, svc *mocks.ShipAssembledProducerService)
	}{
		{
			name:    "успешная сборка заказа",
			wantErr: false,
			args:    event,
			setupMock: func(sleeper *mocks.Sleeper, svc *mocks.ShipAssembledProducerService) {
				sleeper.EXPECT().
					Sleep().
					Return(buildTimeSec)

				svc.EXPECT().
					Produce(ctx, mock.MatchedBy(func(event model.ShipAssembledEvent) bool {
						return event.EventUUID() == eventUUID &&
							event.OrderUUID() == orderUUID &&
							event.UserUUID() == userUUID &&
							event.BuildTimeSec() == buildTimeSec &&
							!event.AssembledAt().IsZero()
					})).
					Return(nil)
			},
		},
		{
			name:    "ошибка продюсера",
			wantErr: true,
			args:    event,
			setupMock: func(sleeper *mocks.Sleeper, svc *mocks.ShipAssembledProducerService) {
				sleeper.EXPECT().
					Sleep().
					Return(buildTimeSec)

				svc.EXPECT().
					Produce(ctx, mock.MatchedBy(func(event model.ShipAssembledEvent) bool {
						return event.EventUUID() == eventUUID &&
							event.OrderUUID() == orderUUID &&
							event.UserUUID() == userUUID &&
							event.BuildTimeSec() == buildTimeSec &&
							!event.AssembledAt().IsZero()
					})).
					Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sleeper := mocks.NewSleeper(t)
			producerSvc := mocks.NewShipAssembledProducerService(t)

			tc.setupMock(sleeper, producerSvc)

			svc := assemblysvc.New(producerSvc, sleeper)

			shipAssembleErr := svc.ShipAssemble(ctx, tc.args)
			if tc.wantErr {
				require.Error(t, shipAssembleErr)
			} else {
				require.NoError(t, shipAssembleErr)
			}
		})
	}
}
