package converter

import (
	"strings"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/input"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func PartProtoToModel(protoPart *inventoryv1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, nil
	}

	var partType model.PartType
	switch model.PartType(strings.TrimPrefix(protoPart.PartType.String(), "PART_TYPE_")) {
	case model.PartTypeWeapon:
		partType = model.PartTypeWeapon
	case model.PartTypeEngine:
		partType = model.PartTypeEngine
	case model.PartTypeHull:
		partType = model.PartTypeHull
	case model.PartTypeShield:
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

func ModelPartsToValidateCompatibilityRequest(uuids input.CreateOrderRequest) *inventoryv1.ValidateCompatibilityRequest {
	pReq := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   uuids.HullUUID,
		EngineUuid: uuids.EngineUUID,
	}

	if uuids.ShieldUUID != nil {
		pReq.ShieldUuid = *uuids.ShieldUUID
	}
	if uuids.WeaponUUID != nil {
		pReq.WeaponUuid = *uuids.WeaponUUID
	}

	return pReq
}
