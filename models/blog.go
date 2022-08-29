package models

import "time"

type Blog struct {
    Id         uint64
    ShopId     uint64
    UserId     uint64
    Title      string
    Images     string
    Content    string
    Liked      uint
    Comments   uint
    IsLiked    uint8     `gorm:"-"`
    LikeList   []string  `gorm:"-"`
    CreateTime time.Time `gorm:"autoCreateTime"`
    UpdateTime time.Time `gorm:"autoUpdateTIme"`
}

func (receiver Blog) TableName() string {
    return "tb_blog"
}
