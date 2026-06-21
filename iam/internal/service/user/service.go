package user

type service struct {
	userRepo  UserRepository
	cryptCost int
}

func New(
	userRepo UserRepository,
	cryptCost int,
) *service {
	return &service{
		userRepo:  userRepo,
		cryptCost: cryptCost,
	}
}
