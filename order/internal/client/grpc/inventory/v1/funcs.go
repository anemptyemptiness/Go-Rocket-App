package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, req model.ListPartsClientRequest) (model.ListPartsClientResponse, error) {
	uuidsStr := make([]string, 0, len(req.UUIDs))
	for _, uuid := range req.UUIDs {
		uuidsStr = append(uuidsStr, uuid.String())
	}

	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuidsStr,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return model.ListPartsClientResponse{}, errs.NewExternalErrWithDescription(errs.ErrInventoryClientNotFound, err.Error())
		case codes.InvalidArgument:
			return model.ListPartsClientResponse{}, errs.NewExternalErrWithDescription(errs.ErrInventoryClientInvalidArgument, err.Error())
		default:
			return model.ListPartsClientResponse{}, errs.NewExternalErrWithDescription(errs.ErrInventoryClientInternal, err.Error())
		}
	}

	return clientConverter.ListPartsClientResponseProtoToModel(resp)
}
