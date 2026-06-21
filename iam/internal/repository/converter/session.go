package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/redis_view"
)

func SessionRedisViewToModel(session redis_view.SessionRedisView) (model.Session, error) {
	modelSession := model.Session{
		UUID:     uuid.MustParse(session.Uuid),
		UserUUID: uuid.MustParse(session.UserUuid),
		Login:    session.Login,
	}

	createdAt, err := time.Parse(time.RFC3339, session.CreatedAt)
	if err != nil {
		return modelSession, err
	}
	modelSession.CreatedAt = createdAt

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return modelSession, err
	}
	modelSession.ExpiresAt = expiresAt

	return modelSession, nil
}
