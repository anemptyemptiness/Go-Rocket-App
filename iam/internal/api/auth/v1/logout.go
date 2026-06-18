package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/converter"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (a *api) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	sessionUUID, err := converter.LogoutRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	err = a.authService.Logout(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	return &authv1.LogoutResponse{}, nil
}
