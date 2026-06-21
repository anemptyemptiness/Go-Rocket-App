package user

import (
	"context"
	"errors"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	userUuid, err := uuid.Parse(userUUID)
	if err != nil {
		return model.User{}, pkgerr.InvalidArgument(errs.ErrInvalidUUID)
	}

	user, err := s.userRepo.GetByUUID(ctx, userUuid)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return model.User{}, pkgerr.NotFound(errs.ErrUserNotFound)
		}
		return model.User{}, pkgerr.Internal(err)
	}

	return user, nil
}
