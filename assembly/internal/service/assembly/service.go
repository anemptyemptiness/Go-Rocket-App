package assembly

type service struct {
	shipAssembledProducerSvc ShipAssembledProducerService
	sleeper                  Sleeper
}

type sleeper struct {
	limit int64
}

func NewSleeper(limit int64) Sleeper {
	return &sleeper{
		limit: limit,
	}
}

func New(
	shipAssembledProducerSvc ShipAssembledProducerService,
	sleeper Sleeper,
) *service {
	return &service{
		shipAssembledProducerSvc: shipAssembledProducerSvc,
		sleeper:                  sleeper,
	}
}
