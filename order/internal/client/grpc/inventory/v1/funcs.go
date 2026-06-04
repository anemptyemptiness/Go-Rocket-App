package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return nil, pkgerr.InvalidArgument(err)
		case codes.Internal:
			return nil, pkgerr.Internal(err)
		default:
			return nil, pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
		}
	}

	return clientConverter.ListPartsClientResponseProtoToModel(resp)
}
