package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func PartProtoToModel(protoPart *inventoryv1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, nil
	}

	partUUID, err := uuid.Parse(protoPart.GetUuid())
	if err != nil {
		return nil, err
	}

	return &model.Part{
		UUID:          partUUID,
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		PartType:      model.PartType(protoPart.GetPartType()),
		StockQuantity: protoPart.GetStockQuantity(),
		CreatedAt:     protoPart.CreatedAt.AsTime(),
	}, nil
}

func ListPartsClientResponseProtoToModel(resp *inventoryv1.ListPartsResponse) (model.ListPartsClientResponse, error) {
	if resp == nil {
		return model.ListPartsClientResponse{}, errs.ErrEmptyResponse
	}

	modelParts := make([]model.Part, 0, len(resp.GetParts()))

	for _, protoPart := range resp.GetParts() {
		modelPart, err := PartProtoToModel(protoPart)
		if err != nil {
			return model.ListPartsClientResponse{}, err
		}
		if modelPart == nil {
			continue
		}

		modelParts = append(modelParts, *modelPart)
	}

	return model.ListPartsClientResponse{
		Parts: modelParts,
	}, nil
}
