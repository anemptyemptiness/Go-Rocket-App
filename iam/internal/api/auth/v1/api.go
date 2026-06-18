package v1

import authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"

type api struct {
	authv1.UnimplementedAuthServiceServer

	authService AuthService
}

func New(
	authService AuthService,
) *api {
	return &api{
		authService: authService,
	}
}
