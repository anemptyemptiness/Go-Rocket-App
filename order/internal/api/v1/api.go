package v1

import orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler

	orderService OrderService
}

func NewAPI(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
