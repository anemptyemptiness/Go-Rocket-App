package auth

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Logout(ctx context.Context, sessionUUID string) error {
	if sessionUUID == "" {
		return pkgerr.InvalidArgument(errs.ErrEmptySessionID)
	}

	sessionUuid, err := uuid.Parse(sessionUUID)
	if err != nil {
		return pkgerr.InvalidArgument(errs.ErrInvalidUUID)
	}

	err = s.sessionRepo.Delete(ctx, sessionUuid)
	if err != nil {
		return pkgerr.Internal(err)
	}

	return nil
}
