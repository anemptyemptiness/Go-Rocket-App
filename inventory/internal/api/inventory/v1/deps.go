package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
