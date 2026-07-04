package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	if sessionUUID == "" {
		slog.ErrorContext(ctx, "whoami", slog.String("error", errs.ErrEmptySessionID.Error()))
		return model.Session{}, model.User{}, pkgerr.InvalidArgument(errs.ErrEmptySessionID)
	}

	sessionUuid, err := uuid.Parse(sessionUUID)
	if err != nil {
		slog.ErrorContext(ctx, "whoami", slog.String("error", errs.ErrInvalidUUID.Error()))
		return model.Session{}, model.User{}, pkgerr.InvalidArgument(errs.ErrInvalidUUID)
	}

	session, err := s.sessionRepo.GetByUUID(ctx, sessionUuid)
	if err != nil {
		slog.ErrorContext(ctx, "whoami", slog.String("error", err.Error()))

		if errors.Is(err, errs.ErrSessionNotFound) {
			return model.Session{}, model.User{}, pkgerr.Unauthenticated(errs.ErrSessionNotFound)
		}
		return model.Session{}, model.User{}, pkgerr.Internal(err)
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		slog.ErrorContext(ctx, "whoami", slog.String("error", err.Error()))

		if errors.Is(err, errs.ErrUserNotFound) {
			return model.Session{}, model.User{}, pkgerr.NotFound(errs.ErrUserNotFound)
		}
		return model.Session{}, model.User{}, pkgerr.Internal(err)
	}

	return session, user, nil
}
