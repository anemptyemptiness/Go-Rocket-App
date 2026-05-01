package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PaymentMethodFromStringToInt32(paymentMethod model.PaymentMethod) int32 {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return int32(paymentv1.PaymentMethod_PAYMENT_METHOD_CARD)
	case model.PaymentMethodSBP:
		return int32(paymentv1.PaymentMethod_PAYMENT_METHOD_SBP)
	case model.PaymentMethodCreditCard:
		return int32(paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD)
	case model.PaymentMethodInvestorMoney:
		return int32(paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY)
	default:
		return int32(paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED)
	}
}
