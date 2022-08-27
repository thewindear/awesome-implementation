package models

import (
	"gorm.io/gorm"
	"implementation-scheme/utils"
	"time"
)

type VoucherOrder struct {
	Id         uint64
	UserId     uint64
	VoucherId  uint64
	PayType    uint8
	Status     uint8
	CreateTime time.Time `gorm:"autoCreateTime"`
	PayTime    time.Time `gorm:"autoCreateTime"`
	UseTime    time.Time `gorm:"type:timestamp"`
	RefundTime time.Time `gorm:"type:timestamp"`
	UpdateTime time.Time `gorm:"autoUpdateTime"`
}

func (receiver *VoucherOrder) BeforeCreate(tx *gorm.DB) (err error) {
	receiver.UseTime = utils.ZeroTime()
	receiver.RefundTime = utils.ZeroTime()
	return nil
}

func (receiver *VoucherOrder) TableName() string {
	return "tb_voucher_order"
}
