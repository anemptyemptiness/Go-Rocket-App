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
		partUUID  = gofakeit.UUID()
		status    = model.OrderStatusPendingPayment
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{
					PartUuid: partUUID,
				},
			},
			Status: status,
		}, nil)

	inventoryClient.EXPECT().
		ReleaseParts(mock.Anything, []string{partUUID}).
		Return(nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusCancelled
	})).Return(nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

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
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{}, errs.ErrOrderNotFound)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

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
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: status,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

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
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: status,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderAlreadyCancelled)
}

func TestCancel_ReleasePartsError(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		partUUID      = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{
					PartUuid: partUUID,
				},
			},
			Status: model.OrderStatusPendingPayment,
		}, nil)

	inventoryClient.EXPECT().
		ReleaseParts(mock.Anything, []string{partUUID}).
		Return(unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
}

func TestCancel_UpdateError(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		partUUID      = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{
					PartUuid: partUUID,
				},
			},
			Status: model.OrderStatusPendingPayment,
		}, nil)

	inventoryClient.EXPECT().
		ReleaseParts(mock.Anything, []string{partUUID}).
		Return(nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusCancelled
	})).Return(unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	err := svc.Cancel(ctx, orderUUID)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
}
