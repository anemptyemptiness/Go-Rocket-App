package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/converter"
	inventoryErrors "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := converter.GetPartRequestProtoToModel(req)
	if err != nil {
		switch {
		case errors.Is(err, inventoryErrors.ErrEmptyRequest):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrPartUUIDIsEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrIncorrectPartUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}

	part, err := a.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		switch {
		case errors.Is(err, inventoryErrors.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartModelToProto(part),
	}, nil
}

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	modelReq, err := converter.ListPartsRequestProtoToModel(req)
	if err != nil {
		switch {
		case errors.Is(err, inventoryErrors.ErrIncorrectPartUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, inventoryErrors.ErrEmptyRequest):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}

	parts, err := a.inventoryService.ListParts(ctx, modelReq)
	if err != nil {
		switch {
		case errors.Is(err, inventoryErrors.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "внутренняя ошибка: %v", err)
		}
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsModelToProto(parts),
	}, nil
}
