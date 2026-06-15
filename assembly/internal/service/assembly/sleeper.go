package assembly

import (
	"math/rand/v2"
	"time"
)

func (s *sleeper) Sleep() int64 {
	sleepTime := rand.Int64N(s.limit + 1)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	return sleepTime
}
