package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
)

type AuthService interface {
	Login(ctx context.Context, loginInput input.LoginInput) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
	Logout(ctx context.Context, sessionUUID string) error
}
