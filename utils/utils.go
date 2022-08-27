package utils

import (
    "math/rand"
    "time"
)

func RandId() uint64 {
    return uint64(rand.Int63n(10000000))
}

func RandIdInt() int {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Int()
}

func ZeroTime() time.Time {
    return time.Date(1970, 12, 31, 0, 0, 0, 0, time.Local)
}
