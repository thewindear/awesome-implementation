package service

import (
    "context"
    "github.com/go-redis/redis/v9"
    "sync"
    "time"
)

// ILock 基于redis实现的锁
type ILock interface {
    Lock(ctx context.Context, key string, val interface{}, expire int) error
    Unlock(ctx context.Context, key string) error
}

type RedisLock struct {
    KeyPrefix string
    rdb       *redis.Client
}

func (l *RedisLock) Lock(ctx context.Context, key string, val interface{}, expire int) (bool, error) {
    return rdb.SetNX(ctx, l.KeyPrefix+key, val, time.Second*time.Duration(expire)).Result()
}

func (l *RedisLock) UnLock(ctx context.Context, key string) (int64, error) {
    return rdb.Del(ctx, l.KeyPrefix+key).Result()
}

// MapLock 基于 sync.Map 实现的锁
type MapLock struct {
    KeyPrefix string
    lock      sync.Map
}

func (l *MapLock) Lock(ctx context.Context, key string, val interface{}, expire int) (bool, error) {
    lockKey := l.KeyPrefix + key
    _, ok := l.lock.Load(lockKey)
    if !ok {
        l.lock.Store(l.KeyPrefix+key, val)
        return true, nil
    }
    return false, nil
}

func (l *MapLock) UnLock(ctx context.Context, key string) (int64, error) {
    l.lock.Delete(l.KeyPrefix + key)
    return 1, nil
}
