package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (c *client) Whoami(ctx context.Context, sessionUUID string) (string, error) {
	resp, err := c.iam.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return "", pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return "", pkgerr.InvalidArgument(err)
		case codes.AlreadyExists:
			return "", pkgerr.AlreadyExists(err)
		case codes.Unauthenticated:
			return "", pkgerr.Unauthenticated(err)
		case codes.Internal:
			return "", pkgerr.Internal(err)
		default:
			return "", pkgerr.Internal(fmt.Errorf("получить информацию о сессии: %w", err))
		}
	}

	if resp != nil && resp.User != nil {
		return resp.User.Uuid, nil
	}
	return "", pkgerr.Internal(errs.ErrExpiredSession)
}
