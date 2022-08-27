package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "implementation-scheme/models"
)

//var couponLock sync.Map
var tryLock = &RedisLock{KeyPrefix: "biz_lock:", rdb: rdb}

type ICouponService interface {
    // 秒杀下单方法
    secKillCoupon(couponId int64)
}

// CouponService 秒杀下单
type CouponService struct {
    db  *gorm.DB
    rdb *redis.Client
}

func (s CouponService) secKillCoupon(ctx context.Context, couponId uint64, userId uint64) error {
    
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
    
    //5. 判断是否已经下过单
    //5.1 判断用户否购买过
    var count int64
    err := s.db.Model(&models.VoucherOrder{}).Where("user_id = ?", userId).Count(&count).Error
    if err != nil {
        return fmt.Errorf("查询是否购买记录失败: %s", err)
    }
    if count > 0 {
        return errors.New("已经购买过无法再次购买")
    }
    /*5.2 使用单机互斥锁 sync.Map 实现
      //5.2 因为这里如果加锁进行判断那么当并发的时候会出现一个人可以下多次单
          _, ok := couponLock.Load(userId)
          if ok {
              return fmt.Errorf("请匆重复操作")
          }
      couponLock.Store(userId, 1)
      defer couponLock.Delete(userId)
      5.2 单机互斥锁 end
    */
    // 5.2 使用redis分布式锁来实现加锁
    var lockVal = fmt.Sprintf("%d:%d", couponId, userId)
    var lockKey = fmt.Sprintf("order:%s", lockVal)
    // lockValue 一般为线程id但是go语言可以为goroutine id
    ok, err := tryLock.ReentryLock(ctx, lockKey, lockVal, 3)
    defer tryLock.ReentryUnlock(ctx, lockKey, lockVal)
    if err != nil {
        return fmt.Errorf("操作失败: %s", err)
    }
    if !ok {
        return errors.New("请匆重复下单")
    }
    err = s.db.Transaction(func(tx *gorm.DB) error {
        //6.扣减库存
        //6.1 这里在高并发情况下会出现超卖的情况
        result = tx.Model(&secKillCoupon).
            Where("stock > 0").
            Update("stock", gorm.Expr("stock - ?", 1))
        if result.Error != nil || result.RowsAffected == 0 {
            return fmt.Errorf("扣减库失败 %s", result.Error)
        }
        //7.创建订单
        order := &models.VoucherOrder{
            UserId:    userId,
            VoucherId: couponId,
        }
        if err = tx.Create(order).Error; err != nil {
            return fmt.Errorf("创建订单失败: %s", err)
        }
        return nil
    })
    return err
}
