package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
	"implementation-scheme/models"
)

type ICouponService interface {
	// 秒杀下单方法
	secKillCoupon(couponId int64)
}

// CouponService 秒杀下单
type CouponService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (s CouponService) secKillCoupon(couponId int64) error {
	var secKillCoupon models.SecKillVoucher
	//1.查询卷是否存在
	result := s.db.Model(&secKillCoupon).Where("voucher_id = ?", couponId).First(&secKillCoupon)
	if result.Error != nil {
		return result.Error
	}
	//2.判断是否开始
	if !secKillCoupon.IsBegin() {
		return errors.New("秒杀没开始")
	}
	//3.判断是否结束
	if secKillCoupon.IsEnd() {
		return errors.New("秒杀已结束")
	}
	//4. 判断库存是否足够
	if !secKillCoupon.StockIsOk() {
		return errors.New("库存不足")
	}
	//5.扣减库存
	//5.1 这里在高并发情况下会出现超卖的情况
	result = s.db.Model(&secKillCoupon).
		Where("stock > 0").
		Update("stock", gorm.Expr("stock - ?", 1))
	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("扣减库失败 %s", result.Error)
	}
	//6.创建订单
	return nil

}
