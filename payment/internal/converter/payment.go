package converter

import (
	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PayOrderRequestProtoToModel(req *paymentv1.PayOrderRequest) (model.PayOrderRequest, error) {
	if req == nil {
		return model.PayOrderRequest{}, errs.ErrEmptyRequest
	}

	return model.PayOrderRequest{
		OrderUUID:     req.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}, nil
}
