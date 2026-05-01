package v1

import (
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func New(inventoryClient inventoryv1.InventoryServiceClient) *client {
	return &client{
		inventoryClient: inventoryClient,
	}
}
