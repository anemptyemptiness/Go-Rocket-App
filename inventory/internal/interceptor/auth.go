package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/client/grpc/iam/v1"
	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/auth"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

const SessionMetadataKey = "session-uuid"

func UnaryAuthInterceptor(iamSvc v1.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, exists := metadata.FromIncomingContext(ctx)
		if !exists {
			return nil, pkgerr.Unauthenticated(errs.ErrEmptyMetadata)
		}

		sessionUUIDs := md.Get(SessionMetadataKey)
		if sessionUUIDs == nil || len(sessionUUIDs) > 1 || sessionUUIDs[0] == "" {
			return nil, pkgerr.Unauthenticated(errs.ErrEmptySessionUUID)
		}

		userUUID, err := iamSvc.Whoami(ctx, sessionUUIDs[0])
		if err != nil {
			return nil, pkgerr.Unauthenticated(errs.ErrExpiredSession)
		}

		userUuid, err := uuid.Parse(userUUID)
		if err != nil {
			return nil, pkgerr.Internal(errs.ErrExpiredSession)
		}

		return handler(auth.WithUserUUID(ctx, userUuid), req)
	}
}
