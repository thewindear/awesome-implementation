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