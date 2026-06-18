package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/config"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/redis_view"
)

func (r *sessionRepository) Create(ctx context.Context, user model.User) (uuid.UUID, error) {
	sessionUUID := uuid.New()

	err := r.db.HSet(
		ctx, redis_view.SessionKey+sessionUUID.String(),
		redis_view.SessionUUIDKey, sessionUUID.String(),
		redis_view.UserUUIDKey, user.UUID.String(),
		redis_view.LoginKey, user.Login,
		redis_view.CreatedAtKey, time.Now().Format(time.RFC3339),
		redis_view.ExpiresAtKey, time.Now().Add(config.AppConfig().Session.TTL).Format(time.RFC3339),
	).Err()
	if err != nil {
		return sessionUUID, err
	}

	err = r.db.Expire(ctx, redis_view.SessionKey+sessionUUID.String(), config.AppConfig().Session.TTL).Err()
	if err != nil {
		return sessionUUID, err
	}

	return sessionUUID, nil
}
