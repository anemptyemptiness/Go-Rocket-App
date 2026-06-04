package part

type service struct {
	inventoryRepo Repository
}

func New(inventoryRepo Repository) *service {
	return &service{
		inventoryRepo: inventoryRepo,
	}
}
