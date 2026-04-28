package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethodString) (uuid.UUID, error) {
	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod(clientConverter.PaymentMethodFromStringToInt32(method)),
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			return uuid.Nil, errs.ErrPaymentClientInvalidArgument
		default:
			return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
		}
	}

	return clientConverter.PayOrderClientResponseProtoToDTO(resp)
}
