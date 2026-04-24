package order

import (
	"context"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/converter"
)

func (r *repository) CreateOrder(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID] = repoConverter.OrderModelToRecord(order)

	return nil
}

func (r *repository) UpdateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[order.OrderUUID]
	if !ok {
		return errs.ErrOrderNotFound
	}

	r.orders[order.OrderUUID] = repoConverter.OrderModelToRecord(order)

	return nil
}
