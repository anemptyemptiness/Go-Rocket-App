package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/converter"
	userv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	userUUID, err := converter.GetUserRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	user, err := a.userService.GetUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: converter.UserModelToProto(user),
	}, nil
}
