package v1

import (
	"google.golang.org/grpc"

	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func New(conn *grpc.ClientConn) *client {
	inventoryClient := inventoryv1.NewInventoryServiceClient(conn)

	return &client{
		inventoryClient: inventoryClient,
	}
}
