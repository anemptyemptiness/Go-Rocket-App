package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PaymentMethodFromModelToProto(paymentMethod model.PaymentMethod) paymentv1.PaymentMethod {
	val, ok := paymentv1.PaymentMethod_value["PAYMENT_METHOD_"+string(paymentMethod)]
	if !ok {
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	return paymentv1.PaymentMethod(val)
}
