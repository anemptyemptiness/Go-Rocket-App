package part

type service struct {
	inventoryRepo Repository
	txManager     TxManager
}

func New(inventoryRepo Repository, txManager TxManager) *service {
	return &service{
		inventoryRepo: inventoryRepo,
		txManager:     txManager,
	}
}
