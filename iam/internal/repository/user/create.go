package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
)

func (r *userRepository) Create(ctx context.Context, login, password string) (uuid.UUID, error) {
	const q = `
		INSERT INTO users (uuid, login, password_hash)
		VALUES ($1, $2, $3);`

	id := uuid.New()

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, q, id, login, password)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return uuid.Nil, errs.ErrUserAlreadyExists
		}
		return uuid.Nil, err
	}

	return id, nil
}
