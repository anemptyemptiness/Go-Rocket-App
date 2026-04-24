package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PayOrderRequestProtoToModel(req *paymentv1.PayOrderRequest) (model.PayOrderRequest, error) {
	if req == nil {
		return model.PayOrderRequest{}, errs.ErrEmptyRequest
	}

	if req.GetOrderUuid() == "" {
		return model.PayOrderRequest{}, errs.ErrOrderUUIDIsEmpty
	}

	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return model.PayOrderRequest{}, errs.ErrPaymentMethodUnspecified
	}

	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return model.PayOrderRequest{}, err
	}

	return model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}, nil
}
