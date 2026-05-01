package converter

import (
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func PartProtoToModel(protoPart *inventoryv1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, nil
	}

	var partType model.PartType
	switch protoPart.GetPartType() {
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		partType = model.PartTypeWeapon
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		partType = model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_HULL:
		partType = model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		partType = model.PartTypeShield
	default:
		partType = model.PartTypeUnspecified
	}

	return &model.Part{
		UUID:          protoPart.GetUuid(),
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		PartType:      partType,
		StockQuantity: protoPart.GetStockQuantity(),
		CreatedAt:     protoPart.CreatedAt.AsTime(),
	}, nil
}

func ListPartsClientResponseProtoToModel(resp *inventoryv1.ListPartsResponse) ([]model.Part, error) {
	if resp == nil {
		return nil, errs.ErrEmptyResponse
	}

	modelParts := make([]model.Part, 0, len(resp.GetParts()))

	for _, protoPart := range resp.GetParts() {
		modelPart, err := PartProtoToModel(protoPart)
		if err != nil {
			return nil, err
		}
		if modelPart == nil {
			continue
		}

		modelParts = append(modelParts, *modelPart)
	}

	return modelParts, nil
}
