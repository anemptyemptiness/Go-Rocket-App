package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/auth"
)

func TestPay_Success(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()
		userUUID        = uuid.New()
		method          = model.PaymentMethodCard
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	items := []model.OrderItem{
		{
			UUID:      gofakeit.UUID(),
			OrderUuid: orderUUID,
			PartUuid:  gofakeit.UUID(),
			PartType:  model.PartTypeShield,
			Price:     10000,
		},
	}

	orderOld := model.Order{
		UUID:       orderUUID,
		Items:      items,
		TotalPrice: 10000,
		UserUUID:   userUUID.String(),
		Status:     model.OrderStatusPendingPayment,
	}

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(orderOld, nil)

	paymentClient.EXPECT().
		PayOrder(mock.Anything, orderUUID, method).
		Return(transactionUUID, nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusPaid &&
			order.TransactionUUID != nil && *order.TransactionUUID == transactionUUID &&
			order.PaymentMethod != nil && *order.PaymentMethod == method
	})).Return(nil)

	producer.EXPECT().
		Produce(ctx, mock.MatchedBy(func(event model.OrderPaidEvent) bool {
			return event.EventUUID() != "" &&
				event.OrderUUID() == orderUUID &&
				event.UserUUID() == userUUID.String()
		})).
		Return(nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.NoError(t, err)
	assert.Equal(t, transactionUUID, trUUID)
}

func TestPay_NotFound(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		userUUID  = uuid.New()
		method    = model.PaymentMethodCard
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{}, errs.ErrOrderNotFound)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderNotFound)
	assert.Empty(t, trUUID)
}

func TestPay_PaymentMethodUnspecified(t *testing.T) {
	t.Parallel()

	var (
		ctx          = context.Background()
		orderUUID    = gofakeit.UUID()
		userUUID     = uuid.New()
		method       = model.PaymentMethodUnspecified
		paymentError = assert.AnError
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}, nil)

	paymentClient.EXPECT().
		PayOrder(mock.Anything, orderUUID, method).
		Return("", paymentError)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, paymentError)
	assert.Empty(t, trUUID)
}

func TestPay_AlreadyPaid(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		userUUID  = uuid.New()
		method    = model.PaymentMethodSBP
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPaid,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderAlreadyPaid)
	assert.Empty(t, trUUID)
}

func TestPay_Cancelled(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
		userUUID  = uuid.New()
		method    = model.PaymentMethodSBP
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusCancelled,
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrOrderAlreadyCancelled)
	assert.Empty(t, trUUID)
}

func TestPay_PaymentServiceError(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		orderUUID     = gofakeit.UUID()
		userUUID      = uuid.New()
		method        = model.PaymentMethodSBP
		unexpectedErr = assert.AnError
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}, nil)

	paymentClient.EXPECT().
		PayOrder(mock.Anything, orderUUID, method).
		Return("", unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
	assert.Empty(t, trUUID)
}

func TestPay_UpdateError(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()
		userUUID        = uuid.New()
		method          = model.PaymentMethodSBP
		unexpectedErr   = assert.AnError
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}, nil)

	paymentClient.EXPECT().
		PayOrder(mock.Anything, orderUUID, method).
		Return(transactionUUID, nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusPaid &&
			order.TransactionUUID != nil && *order.TransactionUUID == transactionUUID &&
			order.PaymentMethod != nil && *order.PaymentMethod == method
	})).Return(unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
	assert.Empty(t, trUUID)
}

func TestPay_ProducerError(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = gofakeit.UUID()
		userUUID        = uuid.New()
		transactionUUID = gofakeit.UUID()
		method          = model.PaymentMethodCard
		unexpectedErr   = assert.AnError
	)

	ctx = auth.WithUserUUID(ctx, userUUID)

	items := []model.OrderItem{
		{
			UUID:      gofakeit.UUID(),
			OrderUuid: orderUUID,
			PartUuid:  gofakeit.UUID(),
			PartType:  model.PartTypeShield,
			Price:     10000,
		},
	}

	orderOld := model.Order{
		UUID:       orderUUID,
		Items:      items,
		TotalPrice: 10000,
		UserUUID:   userUUID.String(),
		Status:     model.OrderStatusPendingPayment,
	}

	repo := mocks.NewOrderRepository(t)
	inventoryClient := mocks.NewInventoryClient(t)
	paymentClient := mocks.NewPaymentClient(t)
	txManager := mocks.NewTxManager(t)
	producer := mocks.NewOrderPaidProducerService(t)

	txManager.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(txCtx context.Context, fn func(ctx context.Context) error) error {
			return fn(txCtx)
		})

	repo.EXPECT().
		GetForUpdate(ctx, orderUUID).
		Return(orderOld, nil)

	paymentClient.EXPECT().
		PayOrder(mock.Anything, orderUUID, method).
		Return(transactionUUID, nil)

	repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderUUID &&
			order.Status == model.OrderStatusPaid &&
			order.TransactionUUID != nil && *order.TransactionUUID == transactionUUID &&
			order.PaymentMethod != nil && *order.PaymentMethod == method
	})).Return(nil)

	producer.EXPECT().
		Produce(ctx, mock.MatchedBy(func(event model.OrderPaidEvent) bool {
			return event.EventUUID() != "" &&
				event.OrderUUID() == orderUUID &&
				event.UserUUID() == userUUID.String()
		})).
		Return(unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient, txManager, producer)

	trUUID, err := svc.Pay(ctx, orderUUID, method)
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
	assert.Empty(t, trUUID)
}
