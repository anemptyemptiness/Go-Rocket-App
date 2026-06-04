package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) GetPart(ctx context.Context, uuidStr string) (model.Part, error) {
	if uuidStr == "" {
		return model.Part{}, pkgerr.InvalidArgument(errs.ErrPartUUIDIsEmpty)
	}

	partUuid, err := uuid.Parse(uuidStr)
	if err != nil || partUuid == uuid.Nil {
		return model.Part{}, pkgerr.InvalidArgument(errs.ErrPartUUIDInvalid)
	}

	part, err := s.inventoryRepo.GetPart(ctx, partUuid.String())
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return model.Part{}, pkgerr.NotFound(err)
		}
		return model.Part{}, pkgerr.Internal(fmt.Errorf("получить деталь: %w", err))
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		for _, uuidCheck := range filter.UUIDs {
			partUuid, err := uuid.Parse(uuidCheck)
			if err != nil || partUuid == uuid.Nil {
				return nil, pkgerr.InvalidArgument(errs.ErrPartUUIDInvalid)
			}
		}
	}

	parts, err := s.inventoryRepo.ListParts(ctx, filter)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, pkgerr.NotFound(err)
		}
		return nil, pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
	}
	if len(filter.UUIDs) > 0 && len(filter.UUIDs) != len(parts) {
		return nil, pkgerr.NotFound(errs.ErrPartNotFound)
	}

	return parts, nil
}
