package utils

import (
	"math/rand"
	"time"
)

func RandId() uint64 {
	return uint64(rand.Int63n(10000000))
}

func ZeroTime() time.Time {
	return time.Date(1970, 12, 31, 0, 0, 0, 0, time.Local)
}
