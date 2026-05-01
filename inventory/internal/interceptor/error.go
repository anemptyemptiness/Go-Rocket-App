package interceptor

import (
	"context"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryErrors "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
)

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		switch {
		case errors.Is(err, inventoryErrors.ErrEmptyRequest):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, inventoryErrors.ErrPartUUIDIsEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrPartUUIDInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}
}
