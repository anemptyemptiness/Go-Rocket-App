package v1

import (
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
)

type client struct {
	iam authv1.AuthServiceClient
}

func New(
	iam authv1.AuthServiceClient,
) *client {
	return &client{
		iam: iam,
	}
}
