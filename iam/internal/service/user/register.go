package user

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

func (s *service) Register(ctx context.Context, registrationInfo input.RegisterInput) (uuid.UUID, error) {
	login := registrationInfo.Info.Info.Login
	password := registrationInfo.Info.Password

	if login == "" {
		slog.ErrorContext(ctx, "register", slog.String("error", errs.ErrInvalidLogin.Error()))
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrInvalidLogin)
	}
	if len([]rune(password)) < 8 {
		slog.ErrorContext(ctx, "register", slog.String("error", errs.ErrWeakPassword.Error()))
		return uuid.Nil, pkgerr.InvalidArgument(errs.ErrWeakPassword)
	}

	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		slog.ErrorContext(ctx, "register", slog.String("error", err.Error()))
		return uuid.Nil, pkgerr.Internal(err)
	}
	if user.Login == login {
		slog.ErrorContext(ctx, "register", slog.String("error", errs.ErrUserAlreadyExists.Error()))
		return uuid.Nil, pkgerr.AlreadyExists(errs.ErrUserAlreadyExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cryptCost)
	if err != nil {
		slog.ErrorContext(ctx, "register", slog.String("error", err.Error()))
		return uuid.Nil, pkgerr.Internal(err)
	}

	userUUID, err := s.userRepo.Create(ctx, login, string(hash))
	if err != nil {
		slog.ErrorContext(ctx, "register", slog.String("error", err.Error()))

		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return uuid.Nil, pkgerr.AlreadyExists(err)
		}
		return uuid.Nil, pkgerr.Internal(err)
	}

	return userUUID, nil
}
