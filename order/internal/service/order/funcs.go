package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return model.Order{}, pkgerr.NotFound(err)
		}
		return model.Order{}, pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
	}

	return order, nil
}

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	if req.HullUUID == "" || req.EngineUUID == "" {
		return model.Order{}, pkgerr.InvalidArgument(errs.ErrHullUUIDAndEngineUUIDAreRequired)
	}
	if len(req.PartUUIDs()) > 0 {
		for _, partUUID := range req.PartUUIDs() {
			if partUUID == "" {
				return model.Order{}, pkgerr.InvalidArgument(errs.ErrInvalidPartUUID)
			}
		}
	}

	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(clientCtx, req.PartUUIDs())
	if err != nil {
		return model.Order{}, err
	}

	var totalPrice int64
	var orderItems []model.OrderItem
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, pkgerr.Conflict(fmt.Errorf("деталь %s: %w", part.Name, errs.ErrPartIsOver))
		}

		orderItems = append(orderItems, model.OrderItem{
			PartUuid: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})

		totalPrice += part.Price
	}

	order := model.Order{
		Items:      orderItems,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
	}

	orderUUID, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return model.Order{}, pkgerr.Internal(fmt.Errorf("создать заказ: %w", err))
	}

	order.SetID(orderUUID)

	return order, err
}

func (s *service) Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return "", pkgerr.NotFound(err)
		}
		return "", pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return "", pkgerr.Conflict(errs.ErrOrderAlreadyPaid)
	case model.OrderStatusCancelled:
		return "", pkgerr.Conflict(errs.ErrOrderAlreadyCancelled)
	case model.OrderStatusAssembled:
		return "", pkgerr.Conflict(errs.ErrOrderAssembled)
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", pkgerr.InvalidArgument(errs.ErrOrderStatusIncorrect)
	}

	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transactionUUID, err := s.paymentClient.PayOrder(clientCtx, orderUUID, method)
	if err != nil {
		return "", fmt.Errorf("оплатить заказ: %w", err)
	}

	order.SetStatus(model.OrderStatusPaid)
	order.SetTransactionID(transactionUUID)
	order.SetPaymentMethod(method)

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return "", pkgerr.NotFound(err)
		}
		return "", pkgerr.Internal(fmt.Errorf("обновить заказ: %w", err))
	}

	return transactionUUID, nil
}

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return pkgerr.NotFound(err)
		}
		return pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return pkgerr.Conflict(errs.ErrOrderAlreadyPaid)
	case model.OrderStatusCancelled:
		return pkgerr.Conflict(errs.ErrOrderAlreadyCancelled)
	case model.OrderStatusAssembled:
		return pkgerr.Conflict(errs.ErrOrderAssembled)
	}

	order.SetStatus(model.OrderStatusCancelled)

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return pkgerr.NotFound(err)
		}
		return pkgerr.Internal(fmt.Errorf("обновить заказ: %w", err))
	}

	return nil
}
