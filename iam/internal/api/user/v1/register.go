package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/converter"
	userv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	registerInput, err := converter.RegisterRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	userUUID, err := a.userService.Register(ctx, registerInput)
	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{
		UserUuid: userUUID.String(),
	}, nil
}
