package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, login, password string) (uuid.UUID, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
