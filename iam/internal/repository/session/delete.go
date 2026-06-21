package session

import (
	"context"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/redis_view"
)

func (r *sessionRepository) Delete(ctx context.Context, sessionUUID uuid.UUID) error {
	err := r.db.Del(ctx, redis_view.SessionKey+sessionUUID.String()).Err()
	if err != nil {
		return err
	}

	return nil
}
