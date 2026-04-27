package middleware

import (
	"net/http"

	"github.com/go-faster/errors"
	ogenmw "github.com/ogen-go/ogen/middleware"

	orderErrors "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func ErrorMiddleware(req ogenmw.Request, next ogenmw.Next) (ogenmw.Response, error) {
	resp, err := next(req)
	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, orderErrors.ErrOrderNotFound):
		return ogenmw.Response{
			Type: &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrShieldUUIDIncorrect):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrWeaponUUIDIncorrect):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrEmptyRequest):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrHullUUIDAndEngineUUIDAreRequired):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrPartIsOver):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrInventoryClientNotFound):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrInventoryClientInvalidArgument):
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrOrderStatusIncorrect):
		return ogenmw.Response{
			Type: &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrPaymentClientInvalidArgument):
		return ogenmw.Response{
			Type: &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrOrderAlreadyPaid):
		return ogenmw.Response{
			Type: &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			},
		}, nil
	case errors.Is(err, orderErrors.ErrOrderAlreadyCancelled):
		return ogenmw.Response{
			Type: &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			},
		}, nil
	default:
		return ogenmw.Response{}, err
	}
}
