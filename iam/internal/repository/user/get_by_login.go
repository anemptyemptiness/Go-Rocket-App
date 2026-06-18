package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errs "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/record"
)

func (r *userRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	const q = `
		SELECT
			u.uuid,
			u.login,
			u.password_hash,
			u.created_at,
			u.updated_at
		FROM users AS u
		WHERE u.login = $1;`

	var user record.User

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, q, login).Scan(
		&user.UUID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("%w: login=%s", errs.ErrUserNotFound, login)
		}
		return model.User{}, err
	}

	return converter.UserRecordToModel(user), nil
}
