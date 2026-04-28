package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

func New() *repository {
	return &repository{
		orders: make(map[uuid.UUID]record.Order),
	}
}
