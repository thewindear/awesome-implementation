package service

import (
    "context"
    "github.com/go-redis/redis/v9"
    "sync"
    "time"
)

var (
    //使用lua完成释放操作
    luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
    //可重入锁释放
    luaReentryUnlock = redis.NewScript(`
local key = KEYS[1];
local threadId = ARGV[1];
local releaseTime = ARGV[2];

-- 判断当前锁是否还是被自己持有
if (redis.call('HEXISTS', key, threadId) == 0) then
    return nil;
end;

-- 是自己的锁 则重入次数-1
local count = redis.call('HINCRBY', key, threadId, -1);
-- 判断是否重入次数为0
if (count > 0) then
    -- 大于0说明还有方法在加锁，重置有效期然后返回
    redis.call('EXPIRE', key, releaseTime)
    return nil;
else -- 所有方法都释放完锁可以删除锁
    redis.call('DEL', key)
    return nil
end;
`)
    //可重入锁lua脚本实现
    luaReentryLock = redis.NewScript(`
local key = KEYS[1]; -- 锁的key
local threadId = ARGV[1]; -- 线程id标识
local releaseTime = ARGV[2]; --释放时间

if (redis.call('exists', key) == 0) then
    -- 不存在设置锁的的线程id
    redis.call('hset', key, threadId, '1');
    -- 设置过期时间
    redis.call('expire', key, releaseTime);
    return 1;
end

-- 锁已经存在时，判断threadId是否为自己
if (redis.call('hexists', key, threadId) == 1) then
    -- 不存在，获取锁，重入次数 +1
    redis.call('hincrby', key, threadId, '1');
    -- 设置有效期
    redis.call('expire', key, releaseTime);
    return 1; -- 返回结果
end
-- 代码走到这里说明获取的锁不是自己的获取锁失败
return 0;
`)
    //唯一性下单
    luaAtomicOrder = redis.NewScript(`
-- 优惠券id
local voucherId = ARGV[1]
-- 用户id
local userId = ARGV[2]

-- 2.数据key
-- 2.1 库存key
local stockKey = 'secKill:stock:' .. voucherId
-- 2.2 订单key 用于保存哪些用户购买了这个券
local orderKey = 'secKill:order:' .. voucherId

-- 3.脚本业务
-- 3.1 判断库存是否充足
if (tonumber(redis.call('GET', stockKey)) <= 0) then
    -- 3.1.1 库存不足 返回1
    return 1
end
-- 3.2 判断用户是否下过单
if (tonumber(redis.call('SISMEMBER', orderKey, userId)) == 1) then
    -- 3.2.1 下过单返回2
    return 2
end
-- 3.3 没下过单 那么扣减库存添加 用户id进指定set集合中
redis.call('INCRBY', stockKey, -1)
redis.call('SADD', orderKey, userId)
-- 下单成功
return 0
`)
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

// ReentryLock 可重入锁
func (l *RedisLock) ReentryLock(ctx context.Context, key string, val interface{}, expire int) (bool, error) {
    res, err := luaReentryLock.Run(ctx, l.rdb, []string{l.KeyPrefix + key}, val, expire).Result()
    if err != nil {
        return false, err
    }
    if res.(int64) == 1 {
        return true, nil
    }
    return false, nil
}

// ReentryUnlock 可重入锁释放锁，先要将计数-1直到为0时才删除key
func (l *RedisLock) ReentryUnlock(ctx context.Context, key string, val interface{}) error {
    _, err := luaReentryUnlock.Run(ctx, l.rdb, []string{l.KeyPrefix + key}, val).Result()
    if err != nil {
        return err
    }
    return nil
}

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
