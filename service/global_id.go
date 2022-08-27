package service

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

// GlobalId 生成全局唯一id的结构体
type GlobalId struct {
	// 生成的id的总位数
	CountBit uint
	// 生成id的开始时间
	BeginTime int64
	// redis
	rdb *redis.Client
}

// NewGlobalId 实例化
func NewGlobalId(countBit uint, rdb *redis.Client) *GlobalId {
	//指定一个开始时间
	beginTime := time.Date(2022, 8, 1, 0, 0, 0, 0, time.Local)
	return &GlobalId{
		CountBit:  countBit,
		BeginTime: beginTime.Unix(),
		rdb:       rdb,
	}
}

// CreateId 生成全局唯一id实现方法
func (g *GlobalId) CreateId(ctx context.Context, keyPrefix string) (int64, error) {
	// 1. 生成时间戳
	timestamp := time.Now().Unix() - g.BeginTime
	// 2. 生成序列号
	// 2.1 这里使用当天日期作为辅助id
	nowDate := time.Now().Format("20060102")
	key := "icr:" + keyPrefix + ":" + nowDate
	count, err := g.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	// 3. 拼接返回
	// 3.1 按位运算进行拼接
	// 3.2 时间戳左移 countBit位 再和总数 或运算
	return timestamp<<g.CountBit | count, err
}
