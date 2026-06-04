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
	partFilter, err := converter.ListPartsRequestProtoToModel(req)
	if err != nil {
		return nil, err
	}

	parts, err := a.inventoryService.ListParts(ctx, partFilter)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsModelToProto(parts),
	}, nil
}

func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (*inventoryv1.ValidateCompatibilityResponse, error) {
	modelReq, err := converter.ValidateCompatibilityRequestToModel(req)
	if err != nil {
		return nil, err
	}

	err = a.inventoryService.ValidateCompatibility(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}

func (a *api) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	modelReq, err := converter.ReservePartsRequestToModel(req)
	if err != nil {
		return nil, err
	}

	err = a.inventoryService.ReserveParts(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}

func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	modelReq, err := converter.ReleasePartsRequestToModel(req)
	if err != nil {
		return nil, err
	}

	err = a.inventoryService.ReleaseParts(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
