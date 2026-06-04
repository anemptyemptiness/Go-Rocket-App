package part

type service struct {
	inventoryRepo        Repository
	compatibilityChecker CompatibilityChecker
	txManager            TxManager
}

func New(
	inventoryRepo Repository,
	compChecker CompatibilityChecker,
	txManager TxManager,
) *service {
	return &service{
		inventoryRepo:        inventoryRepo,
		compatibilityChecker: compChecker,
		txManager:            txManager,
	}
}
