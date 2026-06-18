package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Login(ctx context.Context, loginInput input.LoginInput) (uuid.UUID, error) {
	if loginInput.Login == "" || loginInput.Password == "" {
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrEmptyCredential)
	}

	user, err := s.userRepo.GetByLogin(ctx, loginInput.Login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return uuid.Nil, pkgerr.NotFound(err)
		}
		return uuid.Nil, pkgerr.Internal(err)
	}
	if user.Login != loginInput.Login {
		return uuid.Nil, pkgerr.Unauthenticated(errs.ErrInvalidCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginInput.Password))
	if err != nil {
		return uuid.Nil, pkgerr.Unauthenticated(errs.ErrInvalidCredentials)
	}

	sessionUUID, err := s.sessionRepo.Create(ctx, user)
	if err != nil {
		return uuid.Nil, pkgerr.Internal(err)
	}

	return sessionUUID, nil
}
