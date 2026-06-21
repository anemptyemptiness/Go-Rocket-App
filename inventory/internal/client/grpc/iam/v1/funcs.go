package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (c *client) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	resp, err := c.iam.Whoami(ctx, req)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return nil, pkgerr.InvalidArgument(err)
		case codes.AlreadyExists:
			return nil, pkgerr.AlreadyExists(err)
		case codes.Unauthenticated:
			return nil, pkgerr.Unauthenticated(err)
		case codes.Internal:
			return nil, pkgerr.Internal(err)
		default:
			return nil, pkgerr.Internal(fmt.Errorf("получить информацию о сессии: %w", err))
		}
	}

	return resp, nil
}
