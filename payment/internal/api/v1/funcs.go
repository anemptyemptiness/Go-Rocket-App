package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/converter"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "идентификатор заказа не может быть пустым")
	}

	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "метод оплаты неопределённый")
	}

	modelReq, err := converter.PayOrderRequestProtoToModel(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "конвертация запроса неудачна: %v", err)
	}

	resp, err := a.paymentService.PayOrder(ctx, modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "обработать оплату заказа: %v", err)
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: resp.TransactionUUID.String(),
	}, nil
}
