package models

import "time"

type Follow struct {
    Id           uint64 `gorm:"primaryKey"`
    UserId       uint64
    FollowUserId uint64
    CreateTime   time.Time `gorm:"autoCreateTime"`
}

func (receiver Follow) TableName() string {
    return "tb_follow"
}
