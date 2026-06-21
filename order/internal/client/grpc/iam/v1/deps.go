package v1

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

type Client interface {
	Whoiam(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
}
