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