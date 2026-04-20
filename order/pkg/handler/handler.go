package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// OrderHandler реализует интерфейс orderv1.Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	if req.GetHullUUID() == uuid.Nil || req.GetEngineUUID() == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid и engine_uuid обязательные параметры",
		}, nil
	}

	hullUUID := req.GetHullUUID()
	engineUUID := req.GetEngineUUID()

	uuids := []string{
		hullUUID.String(),
		engineUUID.String(),
	}

	var shieldUUID *uuid.UUID
	if v := req.GetShieldUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "shield_uuid некорректен",
			}, nil
		}
		shieldUUID = &id
		uuids = append(uuids, id.String())
	}

	var weaponUUID *uuid.UUID
	if v := req.GetWeaponUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "weapon_uuid некорректен",
			}, nil
		}
		weaponUUID = &id
		uuids = append(uuids, id.String())
	}

	resp, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case codes.InvalidArgument:
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	var totalPrice int64
	for _, part := range resp.GetParts() {
		if part.GetStockQuantity() <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь " + part.GetName() + " закончилась",
			}, nil
		}
		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.New()
	order := Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
		TotalPrice: totalPrice,
		Status:     string(orderv1.OrderStatusPENDINGPAYMENT),
		CreatedAt:  time.Now(),
	}

	h.store.mu.Lock()
	h.store.orders[orderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay.
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ " + params.OrderUUID.String() + " не найден",
		}, nil
	}

	if order.Status != string(orderv1.OrderStatusPENDINGPAYMENT) {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "некорректный статус заказа",
		}, nil
	}

	paymentMethod := req.GetPaymentMethod()
	var protoPaymentMethod paymentv1.PaymentMethod
	switch paymentMethod {
	case orderv1.PaymentMethodCARD:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "неизвестный метод оплаты",
		}, nil
	}

	resp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: protoPaymentMethod,
	})
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	order.Status = string(orderv1.OrderStatusPAID)

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	paymentMethodStr := string(paymentMethod)
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethodStr

	h.store.mu.Lock()
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel.
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ " + params.OrderUUID.String() + " не найден",
		}, nil
	}

	switch order.Status {
	case string(orderv1.OrderStatusPAID):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	case string(orderv1.OrderStatusCANCELLED):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже отменён",
		}, nil
	}

	h.store.mu.Lock()
	order.Status = string(orderv1.OrderStatusCANCELLED)
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.CancelOrderResponse{}, nil
}
