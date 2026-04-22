package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/converter"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "идентификатор детали не может быть пустым")
	}

	partUUID, err := uuid.Parse(req.GetUuid())
	if err != nil || partUUID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "идентификатор детали невалидный UUID: %v", err)
	}

	part, err := a.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения детали: %v", err)
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartModelToProto(part),
	}, nil
}

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	modelReq, err := converter.ListPartsRequestProtoToModel(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка конвертации запроса: %v", err)
	}

	parts, err := a.inventoryService.ListParts(ctx, modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения списка деталей: %v", err)
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsModelToProto(parts),
	}, nil
}
