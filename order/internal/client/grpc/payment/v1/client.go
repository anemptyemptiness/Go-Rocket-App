package v1

import (
	"google.golang.org/grpc"

	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func NewClient(conn *grpc.ClientConn) *client {
	paymentClient := paymentv1.NewPaymentServiceClient(conn)

	return &client{
		paymentClient: paymentClient,
	}
}
