package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, errs.ErrInventoryClientNotFound
		case codes.InvalidArgument:
			return nil, errs.ErrInventoryClientInvalidArgument
		default:
			return nil, fmt.Errorf("получить список деталей: %w", err)
		}
	}

	return clientConverter.ListPartsClientResponseProtoToModel(resp)
}
