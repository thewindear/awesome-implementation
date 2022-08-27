package service

import (
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
	err := service.secKillCoupon(7)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("下单成功")
	}
}

func TestConcurrentSecKillCoupon(t *testing.T) {
	service := &CouponService{rdb: rdb, db: db}
	wg := sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 200; i++ {
		go func() {
			defer wg.Done()
			err := service.secKillCoupon(7)
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
