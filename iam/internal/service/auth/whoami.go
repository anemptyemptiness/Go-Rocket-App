package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	if sessionUUID == "" {
		return model.Session{}, model.User{}, pkgerr.InvalidArgument(errs.ErrEmptySessionID)
	}

	sessionUuid, err := uuid.Parse(sessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, pkgerr.InvalidArgument(errs.ErrInvalidUUID)
	}

	session, err := s.sessionRepo.GetByUUID(ctx, sessionUuid)
	if err != nil {
		if errors.Is(err, errs.ErrSessionNotFound) {
			return model.Session{}, model.User{}, pkgerr.Unauthenticated(errs.ErrSessionNotFound)
		}
		return model.Session{}, model.User{}, pkgerr.Internal(err)
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return model.Session{}, model.User{}, pkgerr.NotFound(errs.ErrUserNotFound)
		}
		return model.Session{}, model.User{}, pkgerr.Internal(err)
	}

	return session, user, nil
}
