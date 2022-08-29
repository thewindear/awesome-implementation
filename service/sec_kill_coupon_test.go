package service

import (
    "implementation-scheme/models"
    "implementation-scheme/utils"
    "log"
    "testing"
    "time"
)

func TestDB(t *testing.T) {
    row := map[string]interface{}{}
    result := db.Table("tb_voucher").Where("id = ?", 1).Take(&row)
    if result.Error != nil {
        log.Fatalln(result.Error)
    } else {
        t.Log(row)
    }
}

func TestSecKillCoupon(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db}
    err := couponService.secKillCoupon(ctx, 7, utils.RandId())
    if err != nil {
        t.Error(err)
    } else {
        t.Log("下单成功")
    }
}

func TestConcurrentSecKillCoupon(t *testing.T) {
    service := &CouponService{rdb: rdb, db: db}
    ConcurrenceFn(300, func() error {
        return service.secKillCoupon(ctx, 7, utils.RandId())
    })
}

// 新增券优惠券库存
func TestAddCouponStock(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db}
    endTime := time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local)
    err := couponService.AddSecKillCoupon(
        ctx,
        8, "100元代金券", "周一至周日均可使用", "全场通用\\n无需预约\\n可无限叠加\\不兑现、不找零\\n仅限堂食",
        10000, 12000, 1, 1, 500, time.Now(), endTime,
    )
    if err != nil {
        t.Errorf("创建券失败: %s", err)
    } else {
        t.Log("创建券成功")
    }
}

func TestSecKillCouponV2(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db}
    uid := utils.RandIdInt()
    _, err := couponService.secKillCouponV2(ctx, 10, uid)
    if err != nil {
        t.Error(err)
    } else {
        t.Log("下单成功")
    }
}

func TestSecKillCouponV2Concurrent(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db, queue: make(chan *models.VoucherOrder, 1024)}
    ConcurrenceFn(300, func() error {
        uid := utils.RandIdInt()
        order, err := couponService.secKillCouponV2(ctx, 10, uid)
        if err != nil {
            return err
        }
        couponService.queue <- order
        return nil
    })
    couponService.ListenQueue(ctx)
}

func TestSecKillCouponV3Concurrent(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db, queue: make(chan *models.VoucherOrder, 1024)}
    uid := utils.RandIdInt()
    err := couponService.secKillCouponV3(ctx, 10, uid)
    if err != nil {
        t.Error(err)
    } else {
        t.Log("下单成功，等待异步处理")
    }
}

func TestAsyncConsumerOrder(t *testing.T) {
    var couponService = &CouponService{rdb: rdb, db: db, queue: make(chan *models.VoucherOrder, 1024)}
    couponService.ConsumerOrder(ctx)
}