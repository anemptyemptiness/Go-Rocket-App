package converter

import (
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
		PartType:      inventoryv1.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func ListPartsRequestProtoToModel(req *inventoryv1.ListPartsRequest) (model.ListPartsRequest, error) {
	if req == nil {
		return model.ListPartsRequest{}, errs.ErrEmptyRequest
	}

	uuidsStr := make([]string, 0, len(req.GetUuids()))
	uuidsStr = append(uuidsStr, req.GetUuids()...)

	return model.ListPartsRequest{
		UUIDs:    uuidsStr,
		PartType: model.PartType(req.GetPartType()),
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
