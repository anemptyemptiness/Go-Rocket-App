package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/converter"
	orderErrors "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrOrderNotFound):
			return &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.GetOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	return converter.OrderModelToDTO(order), nil
}

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	reqModel, err := converter.CreateOrderRequestToModel(req)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrHullUUIDAndEngineUUIDAreRequired):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrShieldUUIDIncorrect):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrWeaponUUIDIncorrect):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrEmptyRequest):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	respModel, err := a.orderService.Create(ctx, reqModel)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrPartIsOver):
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrInventoryClientNotFound):
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrInventoryClientInvalidArgument):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  respModel.OrderUUID,
		TotalPrice: respModel.TotalPrice,
	}, nil
}

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	payOrderRequestModel, err := converter.PayOrderRequestToModel(req, params)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrEmptyRequest):
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	transactionUUID, err := a.orderService.Pay(ctx, payOrderRequestModel.OrderUUID, payOrderRequestModel.PaymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrOrderStatusIncorrect):
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrOrderNotFound):
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrPaymentClientInvalidArgument):
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrOrderNotFound):
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrOrderAlreadyPaid):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		case errors.Is(err, orderErrors.ErrOrderAlreadyCancelled):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.CancelOrderResponse{}, nil
}
