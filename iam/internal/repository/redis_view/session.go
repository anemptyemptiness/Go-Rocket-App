package redis_view

const (
	SessionKey     = "session:"
	SessionUUIDKey = "uuid"
	UserUUIDKey    = "user_uuid"
	LoginKey       = "login"
	CreatedAtKey   = "created_at"
	ExpiresAtKey   = "expires_at"
)

type SessionRedisView struct {
	Uuid      string `redis:"uuid"`
	UserUuid  string `redis:"user_uuid"`
	Login     string `redis:"login"`
	CreatedAt string `redis:"created_at"`
	ExpiresAt string `redis:"expires_at"`
}
