package user

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func New(
	pool *pgxpool.Pool,
) *userRepository {
	return &userRepository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}
