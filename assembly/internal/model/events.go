package model

import "time"

type OrderPaidEvent interface {
	EventUUID() string
	OrderUUID() string
	UserUUID() string
}

type orderPaidEvent struct {
	eventUUID string
	orderUUID string
	userUUID  string
}

func NewOrderPaidEvent(eventUUID, orderUUID, userUUID string) OrderPaidEvent {
	return &orderPaidEvent{
		eventUUID: eventUUID,
		orderUUID: orderUUID,
		userUUID:  userUUID,
	}
}

func (e *orderPaidEvent) EventUUID() string {
	return e.eventUUID
}

func (e *orderPaidEvent) OrderUUID() string {
	return e.orderUUID
}

func (e *orderPaidEvent) UserUUID() string {
	return e.userUUID
}

type ShipAssembledEvent interface {
	EventUUID() string
	OrderUUID() string
	UserUUID() string
	BuildTimeSec() int64
	AssembledAt() time.Time
	SetBuildTimeSec(secs int64)
	MarkAssembledAt()
}

type shipAssembledEvent struct {
	eventUUID    string
	orderUUID    string
	userUUID     string
	buildTimeSec int64
	assembledAt  time.Time
}

func NewShipAssembledEvent(eventUUID, orderUUID, userUUID string) ShipAssembledEvent {
	return &shipAssembledEvent{
		eventUUID: eventUUID,
		orderUUID: orderUUID,
		userUUID:  userUUID,
	}
}

func (e *shipAssembledEvent) EventUUID() string {
	return e.eventUUID
}

func (e *shipAssembledEvent) OrderUUID() string {
	return e.orderUUID
}

func (e *shipAssembledEvent) UserUUID() string {
	return e.userUUID
}

func (e *shipAssembledEvent) BuildTimeSec() int64 {
	return e.buildTimeSec
}

func (e *shipAssembledEvent) AssembledAt() time.Time {
	return e.assembledAt
}

func (e *shipAssembledEvent) SetBuildTimeSec(secs int64) {
	e.buildTimeSec = secs
}

func (e *shipAssembledEvent) MarkAssembledAt() {
	e.assembledAt = time.Now()
}
