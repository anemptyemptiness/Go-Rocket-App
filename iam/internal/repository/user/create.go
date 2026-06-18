package user

import (
	"context"

	"github.com/google/uuid"
)

func (r *userRepository) Create(ctx context.Context, login, password string) (uuid.UUID, error) {
	const q = `
		INSERT INTO users (uuid, login, password_hash)
		VALUES ($1, $2, $3);`

	id := uuid.New()

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, q, id, login, password)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
