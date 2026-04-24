package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/converter"
)

func (r *repository) GetOrder(_ context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return repoConverter.OrderRecordToModel(order), nil
}
