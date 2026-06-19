package v1

import "context"

type Client interface {
	Whoami(ctx context.Context, sessionUUID string) (string, error)
}
