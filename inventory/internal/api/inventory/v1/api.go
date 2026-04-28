package v1

import inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	inventoryService InventoryService
}

func New(inventoryService InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
