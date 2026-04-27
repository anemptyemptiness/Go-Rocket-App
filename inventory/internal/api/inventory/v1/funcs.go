package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/converter"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := converter.GetPartRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	part, err := a.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartModelToProto(part),
	}, nil
}

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	modelReq, err := converter.ListPartsRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	parts, err := a.inventoryService.ListParts(ctx, modelReq.UUIDs, modelReq.PartType)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsModelToProto(parts),
	}, nil
}
