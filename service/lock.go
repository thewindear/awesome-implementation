package service

import (
    "context"
    "github.com/go-redis/redis/v9"
    "sync"
    "time"
)

var (
    luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
    //使用lua完成释放操作
    luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
    luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

// ILock 基于redis实现的锁
type ILock interface {
    Lock(ctx context.Context, key string, val interface{}, expire int) (bool, error)
    Unlock(ctx context.Context, key string, val interface{})
}

type RedisLock struct {
    KeyPrefix string
    rdb       *redis.Client
}

var _ ILock = &RedisLock{}

func (l *RedisLock) Lock(ctx context.Context, key string, val interface{}, expire int) (bool, error) {
    return rdb.SetNX(ctx, l.KeyPrefix+key, val, time.Second*time.Duration(expire)).Result()
}

func (l *RedisLock) Unlock(ctx context.Context, key string, val interface{}) {
    luaRelease.Run(ctx, l.rdb, []string{l.KeyPrefix + key}, val)
    /**
      keyVal, err := rdb.Get(ctx, l.KeyPrefix+key).Result()
      if err != nil {
          return 0, err
      }
      // 这里需要判断是否为当前线程添加的锁如果不为当前线程添加的锁那么不能做释放
      if value == keyVal {
            //如果在操作del的时候出现阻塞那么就会出现误删除
          return rdb.Del(ctx, l.KeyPrefix+key).Result()
      }
      return 0, nil
    */
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

func (l *MapLock) UnLock(ctx context.Context, key string, val interface{}) (int64, error) {
    l.lock.Delete(l.KeyPrefix + key)
    return 1, nil
}
