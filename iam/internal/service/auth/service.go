package auth

import usersvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user"

type service struct {
	sessionRepo SessionRepository
	userRepo    usersvc.UserRepository
}

func New(
	sessionRepo SessionRepository,
	userRepo usersvc.UserRepository,
) *service {
	return &service{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}
