package interceptor

import "context"

type IAMService interface {
	Whoami(ctx context.Context, sessionUUID string) (string, error)
}
