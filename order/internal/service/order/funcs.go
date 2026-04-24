package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/converter"
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *service) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.CreateOrderResponse, error) {
	uuids := []uuid.UUID{req.HullUUID, req.EngineUUID}

	if req.ShieldUUID != nil {
		uuids = append(uuids, *req.ShieldUUID)
	}
	if req.WeaponUUID != nil {
		uuids = append(uuids, *req.WeaponUUID)
	}

	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.inventoryClient.ListParts(clientCtx, model.ListPartsClientRequest{UUIDs: uuids})
	if err != nil {
		return model.CreateOrderResponse{}, err
	}

	var totalPrice int64
	for _, part := range resp.Parts {
		if part.StockQuantity <= 0 {
			return model.CreateOrderResponse{}, errs.NewExternalErrWithDescription(errs.ErrPartIsOver, part.Name)
		}
		totalPrice += part.Price
	}

	orderUUID := uuid.New()
	err = s.orderRepository.CreateOrder(ctx, model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: req.ShieldUUID,
		WeaponUUID: req.WeaponUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return model.CreateOrderResponse{}, err
	}

	return model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, err
}

func (s *service) PayOrder(ctx context.Context, req model.PayOrderRequest, orderUUID uuid.UUID) (model.PayOrderResponse, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.PayOrderResponse{}, errs.ErrOrderStatusIncorrect
	}

	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.paymentClient.PayOrder(clientCtx, model.PayOrderClientRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: converter.PaymentMethodFromStringToInt32(req.PaymentMethod),
	})
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &resp.TransactionUUID
	order.PaymentMethod = &req.PaymentMethod

	err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	return model.PayOrderResponse(resp), nil
}

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return err
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderAlreadyCancelled
	}

	order.Status = model.OrderStatusCancelled

	err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
