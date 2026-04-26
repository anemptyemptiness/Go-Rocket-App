package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

type Repository interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, uuids []uuid.UUID, partType model.PartType) ([]model.Part, error)
}
