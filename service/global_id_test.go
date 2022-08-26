package service

import (
    "sync"
    "testing"
)

func TestRedisIdWorker(t *testing.T) {
    global := NewGlobalId(32, rdb)
    id, err := global.CreateId(ctx, "order")
    if err != nil {
        t.Error(err)
    } else {
        t.Logf("当前生成的id: %d", id)
    }
}

func TestConcurrentRedisIdWork(t *testing.T) {
    global := NewGlobalId(32, rdb)
    wg := sync.WaitGroup{}
    wg.Add(500)
    for i := 0; i < 500; i++ {
        go func() {
            id, _ := global.CreateId(ctx, "order")
            t.Logf("id: %d", id)
            wg.Done()
        }()
    }
    wg.Wait()
    t.Log("test done")
}
