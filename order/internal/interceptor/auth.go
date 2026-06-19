package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/auth"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

const SessionMetadataKey = "session-uuid"

func SessionUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		sessionUUID, exists := auth.SessionUUIDFromContext(ctx)
		if !exists {
			return pkgerr.Unauthenticated(errs.ErrEmptySessionUUID)
		}

		return invoker(
			metadata.AppendToOutgoingContext(ctx, SessionMetadataKey, sessionUUID),
			method,
			req,
			reply,
			cc,
			opts...,
		)
	}
}
