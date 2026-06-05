package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderservice "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	svcmocks "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order/mocks"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func TestCreate_Success(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		orderUUID  = gofakeit.UUID()
		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
		shieldUUID = gofakeit.UUID()
		weaponUUID = gofakeit.UUID()
		now        = time.Now()
	)

	req := model.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{hullUUID, engineUUID, shieldUUID, weaponUUID}).
		Return([]model.Part{
			{
				UUID:          hullUUID,
				Name:          "hull",
				Description:   "hull-desc",
				Price:         int64(10000),
				PartType:      model.PartTypeHull,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			{
				UUID:          engineUUID,
				Name:          "engine",
				Description:   "engine-desc",
				Price:         int64(20000),
				PartType:      model.PartTypeEngine,
				StockQuantity: 20,
				CreatedAt:     now,
			},
			{
				UUID:          shieldUUID,
				Name:          "shield",
				Description:   "shield-desc",
				Price:         int64(30000),
				PartType:      model.PartTypeShield,
				StockQuantity: 30,
				CreatedAt:     now,
			},
			{
				UUID:          weaponUUID,
				Name:          "weapon",
				Description:   "weapon-desc",
				Price:         int64(40000),
				PartType:      model.PartTypeWeapon,
				StockQuantity: 40,
				CreatedAt:     now,
			},
		}, nil)

	items := []model.OrderItem{
		{PartUuid: hullUUID, Price: 10000, PartType: model.PartTypeHull},
		{PartUuid: engineUUID, Price: 20000, PartType: model.PartTypeEngine},
		{PartUuid: shieldUUID, Price: 30000, PartType: model.PartTypeShield},
		{PartUuid: weaponUUID, Price: 40000, PartType: model.PartTypeWeapon},
	}

	order := model.Order{
		Items:      items,
		TotalPrice: int64(100000),
		Status:     model.OrderStatusPendingPayment,
	}

	repo.EXPECT().
		Create(ctx, order).
		Return(orderUUID, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, orderUUID, order.UUID)
	assert.Equal(t, int64(100000), order.TotalPrice)
	assert.Equal(t, model.OrderStatusPendingPayment, order.Status)
	assert.Len(t, order.Items, len(items))

	for i, item := range order.Items {
		assert.NotEmpty(t, item)
		assert.NotEmpty(t, item.PartUuid)
		assert.Less(t, i, len(items))
		assert.Equal(t, items[i].PartUuid, item.PartUuid)
		assert.Equal(t, items[i].PartType, item.PartType)
		assert.Equal(t, items[i].Price, item.Price)
	}
}

func TestCreate_HullAndEngineRequired(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		hullUUIDEmpty   = ""
		engineUUIDEmpty = ""
	)

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, model.CreateOrderRequest{
		HullUUID:   hullUUIDEmpty,
		EngineUUID: engineUUIDEmpty,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrHullUUIDAndEngineUUIDAreRequired)
	assert.Empty(t, order)
}

func TestCreate_PartNotFound(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
	)

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{hullUUID, engineUUID}).
		Return(nil, pkgerr.NotFound(assert.AnError))

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, model.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, pkgerr.NotFound(assert.AnError))
	assert.Empty(t, order)
}

func TestCreate_PartIsOver(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
	)

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{hullUUID, engineUUID}).
		Return([]model.Part{
			{
				UUID:          hullUUID,
				Name:          "hull",
				Description:   "hull-desc",
				Price:         int64(10000),
				PartType:      model.PartTypeHull,
				StockQuantity: 0,
			},
			{
				UUID:          engineUUID,
				Name:          "engine",
				Description:   "engine-desc",
				Price:         int64(20000),
				PartType:      model.PartTypeEngine,
				StockQuantity: 10,
			},
		}, nil)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, model.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrPartIsOver)
	assert.Empty(t, order)
}

func TestCreate_RepoError(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		hullUUID      = gofakeit.UUID()
		engineUUID    = gofakeit.UUID()
		unexpectedErr = assert.AnError
	)

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{hullUUID, engineUUID}).
		Return([]model.Part{
			{
				UUID:          hullUUID,
				Name:          "hull",
				Description:   "hull-desc",
				Price:         int64(10000),
				PartType:      model.PartTypeHull,
				StockQuantity: 10,
			},
			{
				UUID:          engineUUID,
				Name:          "engine",
				Description:   "engine-desc",
				Price:         int64(20000),
				PartType:      model.PartTypeEngine,
				StockQuantity: 10,
			},
		}, nil)

	order := model.Order{
		Items: []model.OrderItem{
			{
				PartUuid: hullUUID,
				Price:    int64(10000),
				PartType: model.PartTypeHull,
			},
			{
				PartUuid: engineUUID,
				Price:    int64(20000),
				PartType: model.PartTypeEngine,
			},
		},
		TotalPrice: int64(30000),
		Status:     model.OrderStatusPendingPayment,
	}

	repo.EXPECT().
		Create(ctx, order).
		Return("", unexpectedErr)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, model.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, unexpectedErr)
	assert.Empty(t, order)
}

func TestCreate_PartUUIDInvalid(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()
		weaponUUID = ""
	)

	repo := svcmocks.NewOrderRepository(t)
	paymentClient := svcmocks.NewPaymentClient(t)
	inventoryClient := svcmocks.NewInventoryClient(t)

	svc := orderservice.New(repo, paymentClient, inventoryClient)

	order, err := svc.Create(ctx, model.CreateOrderRequest{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		WeaponUUID: &weaponUUID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidPartUUID)
	assert.Empty(t, order)
}
