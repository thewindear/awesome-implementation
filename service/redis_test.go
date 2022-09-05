package service

import (
    "fmt"
    "log"
    "strconv"
    "sync"
    "testing"
)

func TestOneByOneAddRedisKey(t *testing.T) {
    
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("key:%d", i)
        rdb.Set(ctx, key, i, 0)
    }
    fmt.Println("end")
}

func TestBatchAddRedisKey(t *testing.T) {
    
    var groups [10][]string
    for i := 0; i < len(groups); i++ {
        for j := i * 1000; j < i*1000+1000; j++ {
            key := fmt.Sprintf("key:%d", j)
            groups[i] = append(groups[i], key, strconv.Itoa(j))
        }
    }
    var wg sync.WaitGroup
    wg.Add(len(groups))
    for _, group := range groups {
        go func(group []string) {
            pipline := rdb.Pipeline()
            pipline.MSet(ctx, group)
            res, err := pipline.Exec(ctx)
            log.Println(err)
            log.Println(res)
            wg.Done()
        }(group)
    }
    wg.Wait()
}
