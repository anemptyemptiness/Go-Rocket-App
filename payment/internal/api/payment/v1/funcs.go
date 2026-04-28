package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/converter"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	modelReq, err := converter.PayOrderRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, modelReq.OrderUUID, modelReq.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
