package service

import (
    "implementation-scheme/utils"
    "log"
    "sync"
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
    wg := sync.WaitGroup{}
    wg.Add(300)
    for i := 0; i < 300; i++ {
        go func() {
            defer wg.Done()
            err := service.secKillCoupon(ctx, 7, utils.RandId())
            if err != nil {
                t.Logf("下单失败: %s", err)
            } else {
                t.Log("下单成功")
            }
        }()
    }
    wg.Wait()
    t.Log("test done")
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
