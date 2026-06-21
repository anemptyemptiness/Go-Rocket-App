package v1

import authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"

type client struct {
	iamClient authv1.AuthServiceClient
}

func New(
	iamClient authv1.AuthServiceClient,
) *client {
	return &client{
		iamClient: iamClient,
	}
}
