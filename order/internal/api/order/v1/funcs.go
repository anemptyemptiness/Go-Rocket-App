package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/converter"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		return nil, err
	}

	return converter.OrderModelToDTO(order), nil
}

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	reqModel, err := converter.CreateOrderRequestToModel(req)
	if err != nil {
		return nil, err
	}

	order, err := a.orderService.Create(ctx, reqModel)
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  uuid.MustParse(order.UUID),
		TotalPrice: order.TotalPrice,
	}, nil
}

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	payOrderRequestModel, err := converter.PayOrderRequestToModel(req, params)
	if err != nil {
		return nil, err
	}

	transactionUUID, err := a.orderService.Pay(ctx, payOrderRequestModel.OrderUUID, payOrderRequestModel.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID.String())
	if err != nil {
		return nil, err
	}

	return &orderv1.CancelOrderResponse{}, nil
}
