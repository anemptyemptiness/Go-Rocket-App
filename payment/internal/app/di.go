package app

import (
	paymentapi "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/api/payment/v1"
	paymentsvc "github.com/anemptyemptiness/Go-Rocket-App/payment/internal/service/payment"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentAPI     paymentv1.PaymentServiceServer
	paymentService paymentapi.PaymentService
}

func (d *diContainer) PaymentAPI() paymentv1.PaymentServiceServer {
	if d.paymentAPI == nil {
		d.paymentAPI = paymentapi.New(d.PaymentService())
	}

	return d.paymentAPI
}

func (d *diContainer) PaymentService() paymentapi.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentsvc.New()
	}

	return d.paymentService
}
