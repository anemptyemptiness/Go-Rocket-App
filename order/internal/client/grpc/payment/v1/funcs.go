package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: clientConverter.PaymentMethodFromModelToProto(method),
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			return "", pkgerr.InvalidArgument(err)
		default:
			return "", pkgerr.Internal(fmt.Errorf("оплатить заказ: %w", err))
		}
	}

	return resp.GetTransactionUuid(), nil
}
