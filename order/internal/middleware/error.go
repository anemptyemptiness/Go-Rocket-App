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

	status := statusFromError(err)
	if status == http.StatusInternalServerError {
		return ogenmw.Response{}, err
	}

	response, ok := responseFromStatus(req.OperationName, status, err.Error())
	if !ok {
		return ogenmw.Response{}, err
	}

	return response, nil
}

func responseFromStatus(operation string, status int, msg string) (ogenmw.Response, bool) {
	switch operation {
	case "GetOrder":
		return getOrderResponse(status, msg)

	case "CreateOrder":
		return createOrderResponse(status, msg)

	case "PayOrder":
		return payOrderResponse(status, msg)

	case "CancelOrder":
		return cancelOrderResponse(status, msg)

	default:
		return ogenmw.Response{}, false
	}
}

func getOrderResponse(status int, msg string) (ogenmw.Response, bool) {
	switch status {
	case http.StatusNotFound:
		return ogenmw.Response{
			Type: &orderv1.GetOrderNotFound{
				Code:    status,
				Message: msg,
			},
		}, true
	default:
		return ogenmw.Response{}, false
	}
}

func createOrderResponse(status int, msg string) (ogenmw.Response, bool) {
	switch status {
	case http.StatusBadRequest:
		return ogenmw.Response{
			Type: &orderv1.CreateOrderBadRequest{
				Code:    status,
				Message: msg,
			},
		}, true
	case http.StatusNotFound:
		return ogenmw.Response{
			Type: &orderv1.CreateOrderNotFound{
				Code:    status,
				Message: msg,
			},
		}, true
	case http.StatusConflict:
		return ogenmw.Response{
			Type: &orderv1.CreateOrderConflict{
				Code:    status,
				Message: msg,
			},
		}, true
	default:
		return ogenmw.Response{}, false
	}
}

func payOrderResponse(status int, msg string) (ogenmw.Response, bool) {
	switch status {
	case http.StatusBadRequest:
		return ogenmw.Response{
			Type: &orderv1.PayOrderBadRequest{
				Code:    status,
				Message: msg,
			},
		}, true
	case http.StatusConflict:
		return ogenmw.Response{
			Type: &orderv1.PayOrderConflict{
				Code:    status,
				Message: msg,
			},
		}, true
	case http.StatusNotFound:
		return ogenmw.Response{
			Type: &orderv1.PayOrderNotFound{
				Code:    status,
				Message: msg,
			},
		}, true
	default:
		return ogenmw.Response{}, false
	}
}

func cancelOrderResponse(status int, msg string) (ogenmw.Response, bool) {
	switch status {
	case http.StatusConflict:
		return ogenmw.Response{
			Type: &orderv1.CancelOrderConflict{
				Code:    status,
				Message: msg,
			},
		}, true
	case http.StatusNotFound:
		return ogenmw.Response{
			Type: &orderv1.CancelOrderNotFound{
				Code:    status,
				Message: msg,
			},
		}, true
	default:
		return ogenmw.Response{}, false
	}
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, orderErrors.ErrOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, orderErrors.ErrShieldUUIDIncorrect):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrWeaponUUIDIncorrect):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrEmptyRequest):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrHullUUIDAndEngineUUIDAreRequired):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrPartIsOver):
		return http.StatusConflict
	case errors.Is(err, orderErrors.ErrInventoryClientNotFound):
		return http.StatusNotFound
	case errors.Is(err, orderErrors.ErrInventoryClientInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrOrderStatusIncorrect):
		return http.StatusConflict
	case errors.Is(err, orderErrors.ErrPaymentClientInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, orderErrors.ErrOrderAlreadyPaid):
		return http.StatusConflict
	case errors.Is(err, orderErrors.ErrOrderAlreadyCancelled):
		return http.StatusConflict
	case errors.Is(err, orderErrors.ErrInvalidOrderUUID):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
