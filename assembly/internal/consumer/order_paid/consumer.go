package order_paid

import (
	"context"
)

type service struct {
	orderPaidConsumer Consumer
	assembleService   AssembleService
}

func New(
	orderPaidConsumer Consumer,
	assembleService AssembleService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		assembleService:   assembleService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	return s.orderPaidConsumer.Consume(ctx, s.OrderPaid)
}
