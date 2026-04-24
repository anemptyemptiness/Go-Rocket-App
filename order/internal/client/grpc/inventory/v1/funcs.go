package v1

import (
	"context"

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
		return model.ListPartsClientResponse{}, errs.NewInventoryClientInternal(err.Error())
	}

	return clientConverter.ListPartsClientResponseProtoToModel(resp)
}
