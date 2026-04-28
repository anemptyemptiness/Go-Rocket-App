package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ordererrs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		req model.CreateOrderRequest
	}

	type expected struct {
		order   model.Order
		wantErr error
	}

	var (
		ctx           = context.Background()
		hullUUID      = uuid.New()
		engineUUID    = uuid.New()
		shieldUUID    = uuid.New()
		weaponUUID    = uuid.New()
		unexpectedErr = errors.New("внезапность")
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient)
	}{
		{
			name: "успешное создание заказа",
			args: args{
				req: model.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					ShieldUUID: &shieldUUID,
					WeaponUUID: &weaponUUID,
				},
			},
			expected: expected{
				wantErr: nil,
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				inventory.EXPECT().
					ListParts(mock.Anything, []uuid.UUID{hullUUID, engineUUID, shieldUUID, weaponUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Name: "hull", Price: 100000, StockQuantity: 1},
						{UUID: engineUUID, Name: "engine", Price: 200000, StockQuantity: 2},
						{UUID: shieldUUID, Name: "shield", Price: 300000, StockQuantity: 3},
						{UUID: weaponUUID, Name: "weapon", Price: 400000, StockQuantity: 4},
					}, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID != uuid.Nil &&
							order.HullUUID == hullUUID &&
							order.EngineUUID == engineUUID &&
							order.ShieldUUID != nil &&
							*order.ShieldUUID == shieldUUID &&
							order.WeaponUUID != nil &&
							*order.WeaponUUID == weaponUUID &&
							order.TotalPrice == 1000000 &&
							order.Status == model.OrderStatusPendingPayment &&
							!order.CreatedAt.IsZero()
					})).
					Return(nil)
			},
		},
		{
			name: "ошибка: обязательные параметры не переданы",
			args: args{
				req: model.CreateOrderRequest{
					HullUUID:   uuid.Nil,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				order:   model.Order{},
				wantErr: ordererrs.ErrHullUUIDAndEngineUUIDAreRequired,
			},
			setupMock: func(repo *mocks.OrderRepository, payment *mocks.PaymentClient, inventory *mocks.InventoryClient) {},
		},
		{
			name: "ошибка: inventory client not found",
			args: args{
				req: model.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				order:   model.Order{},
				wantErr: ordererrs.ErrInventoryClientNotFound,
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				inventory.EXPECT().
					ListParts(mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return(nil, ordererrs.ErrInventoryClientNotFound)
			},
		},
		{
			name: "ошибка: деталь закончилась",
			args: args{
				req: model.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				order:   model.Order{},
				wantErr: ordererrs.ErrPartIsOver,
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				inventory.EXPECT().
					ListParts(mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Name: "hull", Price: 100000, StockQuantity: 0},
						{UUID: engineUUID, Name: "engine", Price: 200000, StockQuantity: 1},
					}, nil)
			},
		},
		{
			name: "ошибка: создать заказ в репозитории",
			args: args{
				req: model.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				order:   model.Order{},
				wantErr: unexpectedErr,
			},
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient, inventory *mocks.InventoryClient) {
				inventory.EXPECT().
					ListParts(mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Name: "hull", Price: 100000, StockQuantity: 1},
						{UUID: engineUUID, Name: "engine", Price: 200000, StockQuantity: 1},
					}, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(order model.Order) bool {
						return order.OrderUUID != uuid.Nil &&
							order.HullUUID == hullUUID &&
							order.EngineUUID == engineUUID &&
							order.TotalPrice == 300000 &&
							order.Status == model.OrderStatusPendingPayment
					})).
					Return(unexpectedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewOrderRepository(t)
			paymentClient := mocks.NewPaymentClient(t)
			inventoryClient := mocks.NewInventoryClient(t)
			tc.setupMock(repo, paymentClient, inventoryClient)

			svc := ordersvc.New(repo, paymentClient, inventoryClient)

			resp, err := svc.Create(ctx, tc.args.req)
			if tc.expected.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.wantErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp)
				assert.IsType(t, model.Order{}, resp)
				assert.NotEqual(t, uuid.Nil, resp.OrderUUID)
				assert.Equal(t, tc.args.req.HullUUID, resp.HullUUID)
				assert.Equal(t, tc.args.req.EngineUUID, resp.EngineUUID)
				require.NotNil(t, resp.ShieldUUID)
				assert.Equal(t, *tc.args.req.ShieldUUID, *resp.ShieldUUID)
				require.NotNil(t, resp.WeaponUUID)
				assert.Equal(t, *tc.args.req.WeaponUUID, *resp.WeaponUUID)
				assert.Equal(t, int64(1000000), resp.TotalPrice)
				assert.Equal(t, model.OrderStatusPendingPayment, resp.Status)
				assert.False(t, resp.CreatedAt.IsZero())
			}
		})
	}
}
