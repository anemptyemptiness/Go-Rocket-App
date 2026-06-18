package session

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/redis_view"
)

func (r *sessionRepository) GetByUUID(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error) {
	cmd := r.db.HGetAll(ctx, redis_view.SessionKey+sessionUUID.String())

	data, err := cmd.Result()
	if err != nil {
		return model.Session{}, err
	}
	if len(data) == 0 {
		return model.Session{}, errs.ErrSessionNotFound
	}

	var redisSession redis_view.SessionRedisView
	err = cmd.Scan(&redisSession)
	if err != nil {
		return model.Session{}, err
	}

	modelSession, err := converter.SessionRedisViewToModel(redisSession)
	if err != nil {
		return model.Session{}, err
	}

	return modelSession, nil
}
