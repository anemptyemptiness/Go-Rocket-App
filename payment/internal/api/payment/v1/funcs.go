package v1

import (
	"context"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anemptyemptiness/Go-Rocket-App/payment/internal/converter"
	paymentErrors "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	modelReq, err := converter.PayOrderRequestProtoToModel(req)
	if err != nil {
		switch {
		case errors.Is(err, paymentErrors.ErrOrderUUIDIsEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrIncorrectOrderUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrPaymentMethodUnspecified):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrEmptyRequest):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}

	resp, err := a.paymentService.PayOrder(ctx, modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: resp.TransactionUUID.String(),
	}, nil
}
