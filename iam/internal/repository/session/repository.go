package session

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	db  *redis.Client
	ttl time.Duration
}

func New(
	db *redis.Client,
	ttl time.Duration,
) *sessionRepository {
	return &sessionRepository{
		db:  db,
		ttl: ttl,
	}
}
