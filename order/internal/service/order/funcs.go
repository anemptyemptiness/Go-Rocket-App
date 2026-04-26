package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return order, nil
}

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(clientCtx, req.PartUUIDs())
	if err != nil {
		return model.Order{}, fmt.Errorf("получить список деталей: %w", err)
	}

	var totalPrice int64
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrPartIsOver)
		}
		totalPrice += part.Price
	}

	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: req.ShieldUUID,
		WeaponUUID: req.WeaponUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	err = s.orderRepository.Create(ctx, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}

	return order, err
}

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethodString) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return uuid.Nil, errs.ErrOrderStatusIncorrect
	}

	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transactionUUID, err := s.paymentClient.PayOrder(clientCtx, orderUUID, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", err)
	}

	return transactionUUID, nil
}

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderAlreadyCancelled
	}

	order.Status = model.OrderStatusCancelled

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	return nil
}
