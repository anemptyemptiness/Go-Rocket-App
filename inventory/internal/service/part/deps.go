package part

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

type Repository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
