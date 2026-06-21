package v1

import (
	"context"

	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

type Client interface {
	Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error)
}
