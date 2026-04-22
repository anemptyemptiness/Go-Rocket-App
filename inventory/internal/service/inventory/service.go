package inventory

type service struct {
	inventoryRepo Repository
}

func NewService(inventoryRepo Repository) *service {
	return &service{
		inventoryRepo: inventoryRepo,
	}
}
