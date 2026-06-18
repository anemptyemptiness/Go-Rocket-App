package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/converter"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	loginInput, err := converter.LoginRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	sessionUUID, err := a.authService.Login(ctx, loginInput)
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID.String(),
	}, nil
}
