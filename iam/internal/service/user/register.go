package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Register(ctx context.Context, registrationInfo input.RegisterInput) (uuid.UUID, error) {
	login := registrationInfo.Info.Info.Login
	password := registrationInfo.Info.Password

	if login == "" {
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrInvalidLogin)
	}
	if len([]rune(password)) < 8 {
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrWeakPassword)
	}

	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return uuid.Nil, pkgerr.Internal(err)
	}
	if user.Login == login {
		return uuid.Nil, pkgerr.AlreadyExists(errs.ErrUserAlreadyExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cryptCost)
	if err != nil {
		return uuid.Nil, pkgerr.Internal(err)
	}

	userUUID, err := s.userRepo.Create(ctx, login, string(hash))
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return uuid.Nil, pkgerr.AlreadyExists(err)
		}
		return uuid.Nil, pkgerr.Internal(err)
	}

	return userUUID, nil
}
