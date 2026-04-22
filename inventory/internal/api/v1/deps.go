package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, req model.ListPartsRequest) ([]model.Part, error)
}
