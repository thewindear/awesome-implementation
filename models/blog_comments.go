package models

import "time"

type BlogComment struct {
    Id         uint64
    ShopId     uint64
    UserId     uint64
    ParentId   uint64
    AnswerId   uint64
    Content    string
    Liked      uint
    Status     uint8
    CreateTime time.Time `gorm:"autoCreateTime"`
    UpdateTime time.Time `gorm:"autoUpdateTIme"`
}

func (receiver BlogComment) TableName() string {
    return "tb_blog_comments"
}
