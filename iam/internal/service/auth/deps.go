package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, user model.User) (uuid.UUID, error)
	GetByUUID(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	Delete(ctx context.Context, sessionUUID uuid.UUID) error
}
