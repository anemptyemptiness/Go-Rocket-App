package errs

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Если это уже gRPC status - пропускаем как есть.
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		// Стандартные context-ошибки.
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}

		// Доменная ошибка - маппим Code -> grpc code, message - клиенту.
		if be, ok := AsBusinessError(err); ok {
			return nil, status.Error(be.Code().GRPCCode(), be.ClientMessage())
		}

		// Неклассифицированная ошибка - в Internal, в лог с подробностями.
		logger.Error("unexpected grpc error",
			"method", info.FullMethod,
			"error", err,
		)
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}
}
