package converter

import (
	"slices"

	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func PartModelToProto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType.ToProto(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func ListPartsRequestProtoToModel(req *inventoryv1.ListPartsRequest) (model.PartFilter, error) {
	if req == nil {
		return model.PartFilter{}, errs.ErrEmptyRequest
	}

	uuids := slices.Clone(req.GetUuids())

	var partType model.PartType

	switch req.GetPartType() {
	case inventoryv1.PartType_PART_TYPE_HULL:
		partType = model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		partType = model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		partType = model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		partType = model.PartTypeWeapon
	default:
		partType = model.PartTypeUnspecified
	}

	return model.PartFilter{
		UUIDs:    uuids,
		PartType: partType,
	}, nil
}

func PartsModelToProto(parts []model.Part) []*inventoryv1.Part {
	protoParts := make([]*inventoryv1.Part, 0, len(parts))

	for _, part := range parts {
		protoParts = append(protoParts, PartModelToProto(part))
	}

	return protoParts
}

func GetPartRequestProtoToModel(req *inventoryv1.GetPartRequest) (string, error) {
	if req == nil {
		return "", errs.ErrEmptyRequest
	}

	return req.GetUuid(), nil
}
