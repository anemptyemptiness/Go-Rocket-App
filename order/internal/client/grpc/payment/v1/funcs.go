package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, req model.PayOrderClientRequest) (model.PayOrderClientResponse, error) {
	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     req.OrderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			return model.PayOrderClientResponse{}, errs.NewErrPaymentClientInvalidArgument(err.Error())
		default:
			return model.PayOrderClientResponse{}, errs.NewErrPaymentClientInternal(err.Error())
		}
	}

	return clientConverter.PayOrderClientResponseProtoToModel(resp)
}
