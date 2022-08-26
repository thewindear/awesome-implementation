package service

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})
}
