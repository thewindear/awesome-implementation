package service

import (
    "log"
    "testing"
)

func TestRedisSentinel(t *testing.T) {
    log.Println(rdbCluster.Set(ctx, "username:1", "root", 0).Result())
}
