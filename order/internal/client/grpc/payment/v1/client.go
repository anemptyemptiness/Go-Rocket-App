package v1

import (
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func New(paymentClient paymentv1.PaymentServiceClient) *client {
	return &client{
		paymentClient: paymentClient,
	}
}
