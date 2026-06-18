package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
)

type UserService interface {
	Register(ctx context.Context, registrationInfo input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID string) (model.User, error)
}
