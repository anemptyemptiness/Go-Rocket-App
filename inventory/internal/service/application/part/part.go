package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) GetPart(ctx context.Context, uuidStr string) (entity.Part, error) {
	if uuidStr == "" {
		return entity.Part{}, pkgerr.InvalidArgument(errs.ErrPartUUIDIsEmpty)
	}

	partUuid, err := uuid.Parse(uuidStr)
	if err != nil || partUuid == uuid.Nil {
		return entity.Part{}, pkgerr.InvalidArgument(errs.ErrPartUUIDInvalid)
	}

	part, err := s.inventoryRepo.GetPart(ctx, partUuid.String())
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return entity.Part{}, pkgerr.NotFound(err)
		}
		return entity.Part{}, pkgerr.Internal(fmt.Errorf("получить деталь: %w", err))
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter input.PartFilter) ([]entity.Part, error) {
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
		if errors.Is(err, errs.ErrInvalidProperties) {
			return nil, pkgerr.Internal(err)
		}
		return nil, pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
	}
	if len(filter.UUIDs) > 0 && len(filter.UUIDs) != len(parts) {
		return nil, pkgerr.NotFound(errs.ErrPartNotFound)
	}

	return parts, nil
}

func (s *service) CommitParts(ctx context.Context, req input.CommitPartsRequest) error {
	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		filter := input.PartFilter{UUIDs: req.UUIDs}

		parts, err := s.inventoryRepo.ListForUpdate(ctx, filter)
		if err != nil {
			if errors.Is(err, errs.ErrPartNotFound) {
				return pkgerr.NotFound(err)
			}
			if errors.Is(err, errs.ErrInvalidProperties) {
				return pkgerr.Internal(err)
			}
			return pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
		}
		if len(filter.UUIDs) > 0 && len(filter.UUIDs) != len(parts) {
			return pkgerr.NotFound(errs.ErrPartNotFound)
		}

		for _, part := range parts {
			if part.GetStockQuantity() == 0 || part.GetReserved() == 0 {
				return pkgerr.FailedPrecondition(errs.ErrStockQuantityOrReservedIsEmpty)
			}
		}

		err = s.inventoryRepo.CommitParts(txCtx, req)
		if err != nil {
			if errors.Is(err, errs.ErrPartNotFound) {
				return pkgerr.NotFound(err)
			}
			return pkgerr.Internal(fmt.Errorf("списать детали: %w", err))
		}

		return nil
	})
}
