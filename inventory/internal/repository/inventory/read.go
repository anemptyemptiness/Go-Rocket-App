package inventory

import (
	"context"
	"sort"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/record"
)

func (r *repository) GetPart(_ context.Context, uuid uuid.UUID) (model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return repoConverter.PartRecordToModel(part), nil
}

func (r *repository) ListParts(_ context.Context, req model.ListPartsRequest) ([]model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parts := make([]record.Part, 0, len(req.UUIDs))

	if len(req.UUIDs) > 0 {
		for _, UUID := range req.UUIDs {
			part, ok := r.parts[UUID]
			if !ok {
				return nil, errs.ErrPartNotFound
			}

			parts = append(parts, part)
		}
	} else {
		for _, part := range r.parts {
			if req.PartType == model.PartTypeUnspecified || record.PartType(req.PartType) == part.PartType {
				parts = append(parts, part)
			}
		}

		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Name < parts[j].Name
		})
	}

	return repoConverter.PartsRecordToModel(parts), nil
}
