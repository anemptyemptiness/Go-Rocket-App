package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Login(ctx context.Context, loginInput input.LoginInput) (uuid.UUID, error) {
	if loginInput.Login == "" || loginInput.Password == "" {
		slog.ErrorContext(ctx, "логин", slog.String("error", errs.ErrEmptyCredential.Error()))
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrEmptyCredential)
	}

	user, err := s.userRepo.GetByLogin(ctx, loginInput.Login)
	if err != nil {
		slog.ErrorContext(ctx, "логин", slog.String("error", err.Error()))

		if errors.Is(err, errs.ErrUserNotFound) {
			return uuid.Nil, pkgerr.Unauthenticated(errs.ErrInvalidCredentials)
		}
		return uuid.Nil, pkgerr.Internal(err)
	}
	if user.Login != loginInput.Login {
		slog.ErrorContext(ctx, "логин: креды не совпали", slog.String("error", errs.ErrInvalidCredentials.Error()))
		return uuid.Nil, pkgerr.Unauthenticated(errs.ErrInvalidCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginInput.Password))
	if err != nil {
		slog.ErrorContext(ctx, "логин: креды не совпали", slog.String("error", errs.ErrInvalidCredentials.Error()))
		return uuid.Nil, pkgerr.Unauthenticated(errs.ErrInvalidCredentials)
	}

	sessionUUID, err := s.sessionRepo.Create(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "логин", slog.String("error", err.Error()))
		return uuid.Nil, pkgerr.Internal(err)
	}

	return sessionUUID, nil
}
