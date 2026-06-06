package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/input"
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

func (c *client) ValidateCompatibility(ctx context.Context, uuids input.CreateOrderRequest) error {
	_, err := c.inventoryClient.ValidateCompatibility(ctx, clientConverter.ModelPartsToValidateCompatibilityRequest(uuids))
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return pkgerr.InvalidArgument(err)
		case codes.FailedPrecondition:
			return pkgerr.FailedPrecondition(err)
		case codes.ResourceExhausted:
			return pkgerr.ResourceExhausted(err)
		case codes.Internal:
			return pkgerr.Internal(err)
		default:
			return pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
		}
	}

	return nil
}

func (c *client) ReserveParts(ctx context.Context, uuids []string) error {
	_, err := c.inventoryClient.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{Uuids: uuids})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return pkgerr.InvalidArgument(err)
		case codes.FailedPrecondition:
			return pkgerr.FailedPrecondition(err)
		case codes.ResourceExhausted:
			return pkgerr.ResourceExhausted(err)
		case codes.Internal:
			return pkgerr.Internal(err)
		default:
			return pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
		}
	}

	return nil
}

func (c *client) ReleaseParts(ctx context.Context, uuids []string) error {
	_, err := c.inventoryClient.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{Uuids: uuids})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return pkgerr.NotFound(err)
		case codes.InvalidArgument:
			return pkgerr.InvalidArgument(err)
		case codes.FailedPrecondition:
			return pkgerr.FailedPrecondition(err)
		case codes.ResourceExhausted:
			return pkgerr.ResourceExhausted(err)
		case codes.Internal:
			return pkgerr.Internal(err)
		default:
			return pkgerr.Internal(fmt.Errorf("получить список деталей: %w", err))
		}
	}

	return nil
}
