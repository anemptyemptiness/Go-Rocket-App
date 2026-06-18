package user

type service struct {
	userRepo UserRepository
}

func New(
	userRepo UserRepository,
) *service {
	return &service{
		userRepo: userRepo,
	}
}
