package models

import "time"

type Voucher struct {
    Id          uint64 `gorm:"primaryKey"`
    ShopId      uint64
    Title       string `gorm:"size:255"`
    SubTitle    string `gorm:"size:255"`
    Rules       string `gorm:"size:1024"`
    PayValue    uint64
    ActualValue int64
    Type        uint8
    Status      uint8
    CreateTime  time.Time `gorm:"autoCreateTime"`
    UpdateTime  time.Time `gorm:"autoUpdateTime"`
}

func (receiver Voucher) TableName() string {
    return "tb_voucher"
}
