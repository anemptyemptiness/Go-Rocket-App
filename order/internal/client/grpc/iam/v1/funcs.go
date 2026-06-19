package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/iam/v1/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (c *client) Whoiam(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	resp, err := c.iamClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return model.Session{}, model.User{}, pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return model.Session{}, model.User{}, pkgerr.InvalidArgument(err)
		case codes.AlreadyExists:
			return model.Session{}, model.User{}, pkgerr.AlreadyExists(err)
		case codes.Unauthenticated:
			return model.Session{}, model.User{}, pkgerr.Unauthenticated(err)
		case codes.Internal:
			return model.Session{}, model.User{}, pkgerr.Internal(err)
		default:
			return model.Session{}, model.User{}, pkgerr.Internal(fmt.Errorf("получить информацию о сессии: %w", err))
		}
	}

	return converter.SessionProtoToModel(resp.GetSession()), converter.UserProtoToModel(resp.GetUser()), nil
}
