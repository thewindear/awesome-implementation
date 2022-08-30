package models

type Sign struct {
    Id       uint64 `gorm:"primaryKey"`
    UserId   uint64
    Year     int
    Month    uint8
    Date     uint8
    IsBackup uint8
}

func (receiver Sign) TableName() string {
    return "tb_sign"
}
