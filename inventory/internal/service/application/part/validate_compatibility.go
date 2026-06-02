package part

import (
	"context"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) ValidateCompatibility(ctx context.Context, req input.ValidateCompatibilityRequest) error {
	err := s.validateCompatibility(ctx, req)
	if err != nil {
		return err
	}

	parts, err := s.ListParts(ctx, input.PartFilter{UUIDs: req.UUIDs()})
	if err != nil {
		return err
	}

	err = s.resolveShipSlots(ctx, parts, req)
	if err != nil {
		return err
	}

	err = s.compatibilityChecker.Check(parts)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) validateCompatibility(_ context.Context, req input.ValidateCompatibilityRequest) error {
	slots := map[string]string{
		"hull":   req.HullUUID,
		"engine": req.EngineUUID,
	}

	if req.ShieldUUID != "" {
		slots["shield"] = req.ShieldUUID
	}
	if req.WeaponUUID != "" {
		slots["weapon"] = req.WeaponUUID
	}

	seen := make(map[string]struct{})
	for _, id := range slots {
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			return pkgerr.InvalidArgument(errs.ErrPartTypeMismatch)
		}

		seen[id] = struct{}{}
	}

	return nil
}

func (s *service) resolveShipSlots(_ context.Context, parts []entity.Part, req input.ValidateCompatibilityRequest) error {
	for _, part := range parts {
		switch {
		case part.GetPartUUID() == req.HullUUID && part.GetPartType() != valueobject.PartTypeHull:
			return pkgerr.InvalidArgument(errs.ErrPartTypeMismatch)
		case part.GetPartUUID() == req.EngineUUID && part.GetPartType() != valueobject.PartTypeEngine:
			return pkgerr.InvalidArgument(errs.ErrPartTypeMismatch)
		case part.GetPartUUID() == req.ShieldUUID && part.GetPartType() != valueobject.PartTypeShield:
			return pkgerr.InvalidArgument(errs.ErrPartTypeMismatch)
		case part.GetPartUUID() == req.WeaponUUID && part.GetPartType() != valueobject.PartTypeWeapon:
			return pkgerr.InvalidArgument(errs.ErrPartTypeMismatch)
		}
	}

	return nil
}

func (s *service) ReserveParts(ctx context.Context, req input.ReservePartsRequest) error {
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		parts, err := s.ListParts(txCtx, input.PartFilter{UUIDs: req.UUIDs})
		if err != nil {
			return err
		}

		for idx := range parts {
			if err = parts[idx].Reserve(); err != nil {
				return err
			}
		}

		if err = s.inventoryRepo.UpdateReservedBatch(txCtx, parts); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ReleaseParts(ctx context.Context, req input.ReleasePartsRequest) error {
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		parts, err := s.ListParts(txCtx, input.PartFilter{UUIDs: req.UUIDs})
		if err != nil {
			return err
		}

		for idx := range parts {
			if err = parts[idx].Release(); err != nil {
				return err
			}
		}

		if err = s.inventoryRepo.UpdateReservedBatch(txCtx, parts); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
