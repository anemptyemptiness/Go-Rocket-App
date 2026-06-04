package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{
					UUID:      gofakeit.UUID(),
					OrderUuid: orderUUID,
					PartUuid:  gofakeit.UUID(),
					PartType:  model.PartTypeWeapon,
					Price:     10000,
					CreatedAt: time.Now(),
				},
			},
			TotalPrice: 10000,
			Status:     model.OrderStatusPendingPayment,
			CreatedAt:  time.Now(),
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Get(ctx, orderUUID)
	require.NoError(t, err)
	assert.Equal(t, orderUUID, order.UUID)
	assert.Equal(t, int64(10000), order.TotalPrice)
	assert.Equal(t, model.OrderStatusPendingPayment, order.Status)
	assert.Len(t, order.Items, 1)
	assert.NotEmpty(t, order.Items[0].UUID)
}

func TestGet_NotFound(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{}, errs.ErrOrderNotFound)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Get(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderNotFound)
	assert.Empty(t, order)
}

func TestGet_RepoError(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{}, unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Get(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
	assert.Empty(t, order)
}
