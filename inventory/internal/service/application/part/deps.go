package part

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

type Repository interface {
	GetPart(ctx context.Context, uuid string) (entity.Part, error)
	ListParts(ctx context.Context, filter input.PartFilter) ([]entity.Part, error)
	ListForUpdate(ctx context.Context, filter input.PartFilter) ([]entity.Part, error)
	UpdateReservedBatch(ctx context.Context, parts []entity.Part) error
	CommitParts(ctx context.Context, req input.CommitPartsRequest) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CompatibilityChecker interface {
	Check(parts []entity.Part) error
}
