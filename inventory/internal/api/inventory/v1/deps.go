package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (entity.Part, error)
	ListParts(ctx context.Context, filter input.PartFilter) ([]entity.Part, error)
	ValidateCompatibility(ctx context.Context, req input.ValidateCompatibilityRequest) error
	ReserveParts(ctx context.Context, req input.ReservePartsRequest) error
	ReleaseParts(ctx context.Context, req input.ReleasePartsRequest) error
}
