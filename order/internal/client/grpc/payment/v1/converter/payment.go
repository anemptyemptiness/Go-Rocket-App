package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PayOrderClientResponseProtoToDTO(resp *paymentv1.PayOrderResponse) (uuid.UUID, error) {
	if resp == nil {
		return uuid.Nil, errs.ErrEmptyResponse
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}

func PaymentMethodFromStringToInt32(paymentMethod model.PaymentMethodString) model.PaymentMethodInt32 {
	switch paymentMethod {
	case model.PaymentMethodStringCard:
		return model.PaymentMethodInt32Card
	case model.PaymentMethodStringSBP:
		return model.PaymentMethodInt32SBP
	case model.PaymentMethodStringCreditCard:
		return model.PaymentMethodInt32CreditCard
	case model.PaymentMethodStringInvestorMoney:
		return model.PaymentMethodInt32InvestorMoney
	default:
		return model.PaymentMethodInt32Unspecified
	}
}
