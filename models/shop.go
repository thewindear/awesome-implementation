package models

import "time"

type Shop struct {
    Id         uint64 `gorm:"primaryKey"`
    Name       string `gorm:"size:128"`
    TypeId     uint64
    Images     string `gorm:"size:1024"`
    Area       string `gorm:"size:128"`
    Address    string `gorm:"size:255"`
    X          float64
    Y          float64
    AvgPrice   int
    Score      int
    ShopDist   float64   `gorm:"-"`
    CreateTime time.Time `gorm:"autoCreateTime"`
    UpdateTime time.Time `gorm:"autoUpdateTime"`
}

func (receiver Shop) TableName() string {
    return "tb_shop"
}
