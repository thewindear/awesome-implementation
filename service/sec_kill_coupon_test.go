package service

import (
    "implementation-scheme/utils"
    "log"
    "sync"
    "testing"
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
    service := &CouponService{rdb: rdb, db: db}
    err := service.secKillCoupon(ctx, 7, utils.RandId())
    if err != nil {
        t.Error(err)
    } else {
        t.Log("下单成功")
    }
}

func TestConcurrentSecKillCoupon(t *testing.T) {
    service := &CouponService{rdb: rdb, db: db}
    wg := sync.WaitGroup{}
    wg.Add(100)
    for i := 0; i < 100; i++ {
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

func TestUseRedisLua(t *testing.T) {
    var tryLock2 = &RedisLock{KeyPrefix: "biz_lock:", rdb: rdb}
    ok, err := tryLock2.Lock(ctx, "order:abcdefg222", "abcdefg", 20)
    t.Log(ok)
    t.Log(err)
    tryLock2.Unlock(ctx, "order:abcdefg222", "abcdefg")
    //defer tryLock.Unlock(ctx, "order:abcdefg", "abcdefg")
}
