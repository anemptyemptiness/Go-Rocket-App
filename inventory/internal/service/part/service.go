package part

type service struct {
	inventoryRepo Repository
}

func NewService(inventoryRepo Repository) *service {
	return &service{
		inventoryRepo: inventoryRepo,
	}
}
