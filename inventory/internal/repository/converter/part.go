package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/record"
)

func PartRecordToModel(part record.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}

func PartsRecordToModel(parts []record.Part) []model.Part {
	modelParts := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		modelParts = append(modelParts, PartRecordToModel(part))
	}

	return modelParts
}
