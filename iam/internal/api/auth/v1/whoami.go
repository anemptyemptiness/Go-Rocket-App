package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/converter"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	sessionUUID, err := converter.WhoamiRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	session, user, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	return &authv1.WhoamiResponse{
		Session: converter.SessionModelToProto(session),
		User:    converter.UserModelToProto(user),
	}, nil
}
