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

	partUUID, err := uuid.Parse(uuidStr)
	if err != nil || partUUID == uuid.Nil {
		return model.Part{}, errs.ErrIncorrectPartUUID
	}

	part, err := s.inventoryRepo.GetPart(ctx, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("получать деталь: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, uuidsStr []string, partType model.PartType) ([]model.Part, error) {
	uuids := make([]uuid.UUID, 0, len(uuidsStr))
	for _, uuidStr := range uuidsStr {
		parsedUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, errs.ErrIncorrectPartUUID
		}

		uuids = append(uuids, parsedUUID)
	}

	parts, err := s.inventoryRepo.ListParts(ctx, uuids, partType)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return parts, nil
}
