package models

import (
    "time"
)

// SecKillVoucher 秒杀配置信息
type SecKillVoucher struct {
    VoucherId  uint64 `gorm:"primaryKey"`
    Stock      int64
    CreateTime time.Time `gorm:"autoCreateTime"`
    BeginTime  time.Time
    EndTime    time.Time
    UpdateTime time.Time `gorm:"autoUpdateTime"`
}

func (s SecKillVoucher) StockIsOk() bool {
    return s.Stock > 0
}

// IsBegin 是否开始
func (s SecKillVoucher) IsBegin() bool {
    return time.Now().After(s.BeginTime)
}

// IsEnd 是否结束
func (s SecKillVoucher) IsEnd() bool {
	return time.Now().After(s.EndTime)
}

func (s SecKillVoucher) TableName() string {
	return "tb_seckill_voucher"
}
