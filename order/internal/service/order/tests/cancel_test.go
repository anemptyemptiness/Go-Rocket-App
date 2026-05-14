package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
)

func TestCancel_Success(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		status    = model.OrderStatusPendingPayment
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: status,
		}, nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusCancelled
	})).Return(nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	err := svc.Cancel(ctx, orderUUID)
	require.NoError(t, err)
}

func TestCancel_NotFound(t *testing.T) {
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

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderNotFound)
}

func TestCancel_AlreadyPaid(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		status    = model.OrderStatusPaid
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: status,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderAlreadyPaid)
}

func TestCancel_AlreadyCancelled(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		status    = model.OrderStatusCancelled
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)

	repo.EXPECT().
		Get(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: status,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderAlreadyCancelled)
}
