package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentErrors "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/errors"
)

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		switch {
		case errors.Is(err, paymentErrors.ErrEmptyRequest):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrOrderUUIDIsEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrPaymentMethodUnspecified):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, paymentErrors.ErrIncorrectOrderUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}
}
