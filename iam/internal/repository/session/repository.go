package session

import "github.com/redis/go-redis/v9"

type sessionRepository struct {
	db *redis.Client
}

func New(db *redis.Client) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}
