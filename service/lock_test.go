package service

import (
    "github.com/bsm/redislock"
    "log"
    "testing"
    "time"
)

func TestUseRedisLua(t *testing.T) {
    var tryLock2 = &RedisLock{KeyPrefix: "biz_lock:", rdb: rdb}
    ok, err := tryLock2.Lock(ctx, "order:abcdefg222", "abcdefg", 20)
    t.Log(ok)
    t.Log(err)
    tryLock2.Unlock(ctx, "order:abcdefg222", "abcdefg")
    //defer tryLock.Unlock(ctx, "order:abcdefg", "abcdefg")
}

func TestRedisLock(t *testing.T) {
    locker := redislock.New(rdb)
    lock, err := locker.Obtain(ctx, "my-key", 100*time.Second, nil)
    if err == redislock.ErrNotObtained {
        t.Fatal("无法获取锁")
    } else if err != nil {
        t.Fatal(err)
    }
    t.Log(lock.Metadata())
}

func TestRetryLock(t *testing.T) {
    var retryLock = &RedisLock{KeyPrefix: "biz_lock2:", rdb: rdb}
    ok, err := retryLock.ReentryLock(ctx, "order:20220901", "20220901", 100)
    defer retryLock.ReentryUnlock(ctx, "order:20220901", "20220901")
    if err != nil {
        log.Println(err)
    } else {
        t.Log("第1次加锁", ok)
        retryLock2Fn(retryLock, t)
    }
}

func retryLock2Fn(rdbLock *RedisLock, t *testing.T) {
    ok, err := rdbLock.ReentryLock(ctx, "order:20220901", "20220901", 100)
    defer rdbLock.ReentryUnlock(ctx, "order:20220901", "20220901")
    if err != nil {
        log.Println(err)
    } else {
        t.Log("第2次加锁", ok)
    }
}
