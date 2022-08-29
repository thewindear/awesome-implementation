-- 优惠券id
local voucherId = ARGV[1]
-- 用户id
local userId = ARGV[2]
-- 订单id
-- local orderId = ARGV[3]

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

-- 3.4 发送消息到队列中，xadd stream.orders * k1 v1 k2 v2
-- redis.call('xadd', 'stream.orders', '*', 'userId', userId, 'voucherId', voucherId, 'id', orderId)
redis.call('xadd', 'stream.orders', '*', 'userId', userId, 'voucherId', voucherId)

-- 下单成功
return 0