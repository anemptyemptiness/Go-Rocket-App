package auth

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Logout(ctx context.Context, sessionUUID string) error {
	if sessionUUID == "" {
		slog.ErrorContext(ctx, "лог аут", slog.String("error", errs.ErrEmptySessionID.Error()))
		return pkgerr.InvalidArgument(errs.ErrEmptySessionID)
	}

	sessionUuid, err := uuid.Parse(sessionUUID)
	if err != nil {
		slog.ErrorContext(ctx, "лог аут", slog.String("error", errs.ErrInvalidUUID.Error()))
		return pkgerr.InvalidArgument(errs.ErrInvalidUUID)
	}

	err = s.sessionRepo.Delete(ctx, sessionUuid)
	if err != nil {
		slog.ErrorContext(ctx, "лог аут", slog.String("error", err.Error()))
		return pkgerr.Internal(err)
	}

	return nil
}
