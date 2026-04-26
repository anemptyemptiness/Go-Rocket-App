package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	part, err := s.inventoryRepo.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, fmt.Errorf("получать деталь: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, req model.ListPartsRequest) ([]model.Part, error) {
	parts, err := s.inventoryRepo.ListParts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return parts, nil
}
