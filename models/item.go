package models

import "time"

type Item struct {
    Id         uint64    `gorm:"primaryKey" json:"id"`
    Title      string    `gorm:"size:264" json:"title"`
    Name       string    `gorm:"size:128" json:"name"`
    Price      uint64    `json:"price"`
    Image      string    `gorm:"size:200" json:"image"`
    Category   string    `gorm:"size:200" json:"category"`
    Brand      string    `gorm:"size:100" json:"brand"`
    Spec       string    `gorm:"size:200" json:"spec"`
    Status     uint8     `gorm:"default:1" json:"status"`
    CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
    UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

func (i Item) TableName() string {
    return "tb_item"
}

type ItemStock struct {
    ItemId uint64 `gorm:"primaryKey" json:"itemId"`
    Stock  int    `gorm:"default:9999" json:"stock"`
    Sold   int    `gorm:"default:0" json:"sold"`
}

func (s ItemStock) TableName() string {
    return "tb_item_stock"
}
