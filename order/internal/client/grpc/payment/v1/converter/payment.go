package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func PayOrderClientResponseProtoToModel(resp *paymentv1.PayOrderResponse) (model.PayOrderClientResponse, error) {
	if resp == nil {
		return model.PayOrderClientResponse{}, errs.ErrEmptyResponse
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return model.PayOrderClientResponse{}, err
	}

	return model.PayOrderClientResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
