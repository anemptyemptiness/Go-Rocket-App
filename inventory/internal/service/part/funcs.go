package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuidStr string) (model.Part, error) {
	if uuidStr == "" {
		return model.Part{}, errs.ErrPartUUIDIsEmpty
	}

	partUuid, err := uuid.Parse(uuidStr)
	if err != nil || partUuid == uuid.Nil {
		return model.Part{}, errs.ErrPartUUIDInvalid
	}

	part, err := s.inventoryRepo.GetPart(ctx, uuidStr)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.Uuids) > 0 {
		for _, uuidCheck := range filter.Uuids {
			partUuid, err := uuid.Parse(uuidCheck)
			if err != nil || partUuid == uuid.Nil {
				return nil, errs.ErrPartUUIDInvalid
			}
		}
	}

	parts, err := s.inventoryRepo.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	if len(filter.Uuids) > 0 && len(parts) != len(filter.Uuids) {
		return nil, errs.ErrPartNotFound
	}

	return parts, nil
}
