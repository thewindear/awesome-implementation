package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "implementation-scheme/models"
    "log"
    "strconv"
    "time"
)

//var couponLock sync.Map

type ICouponService interface {
    // 秒杀下单方法
    secKillCoupon(couponId int64)
}

// CouponService 秒杀下单
type CouponService struct {
    db    *gorm.DB
    rdb   *redis.Client
    queue chan *models.VoucherOrder
}

const (
    SecKillStockKey = "seckill:stock:%d"
)

// AddSecKillCoupon 添加优惠券信息
// 并将库存信息保存至redis中
func (s CouponService) AddSecKillCoupon(ctx context.Context,
    shopId uint64, title,
    subString, rules string,
    payValue uint64,
    actualValue int64,
    couponType, status uint8, stock int64, beginTime, endTime time.Time) error {
    //保存券信息
    coupon := &models.Voucher{
        ShopId:      shopId,
        Title:       title,
        SubTitle:    subString,
        Rules:       rules,
        PayValue:    payValue,
        ActualValue: actualValue,
        Type:        couponType,
        Status:      status,
    }
    secKillVoucher := &models.SecKillVoucher{}
    return s.db.Transaction(func(tx *gorm.DB) error {
        err := tx.Create(coupon).Error
        if err != nil {
            return err
        }
        secKillVoucher.VoucherId = coupon.Id
        secKillVoucher.Stock = stock
        secKillVoucher.BeginTime = beginTime
        secKillVoucher.EndTime = endTime
    
        err = tx.Create(secKillVoucher).Error
        if err != nil {
            return err
        }
        //将秒杀库存放到redis中
        stockKey := fmt.Sprintf(SecKillStockKey, coupon.Id)
        s.rdb.Set(ctx, stockKey, secKillVoucher.Stock, -1)
        return nil
    })
}

// ConsumerOrder 获取异步消息然后消费
func (s CouponService) ConsumerOrder(ctx context.Context) {
    // 1. 获取消息中的队列订单信息
    // 2. 判断消息获取是否成功
    // 2.1 如果获取失败谙有没有消息，继续下次循环
    // 3 如果获取成功 下单
    // 4 ack确认
    var count = 0
    s.rdb.XGroupCreate(ctx, "stream.orders", "g1", "0")
    for {
        cmd := s.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    "g1",
            Consumer: "g1-c1",
            Streams:  []string{"stream.orders", ">"},
            Count:    1,
            Block:    time.Second,
            NoAck:    false,
        })
        streamRes, err := cmd.Result()
        if err == redis.Nil {
            if count >= 10 {
                break
            }
            time.Sleep(time.Millisecond * 20)
            count += 1
            continue
        }
        userId, _ := strconv.Atoi(streamRes[0].Messages[0].Values["userId"].(string))
        voucherId, _ := strconv.Atoi(streamRes[0].Messages[0].Values["voucherId"].(string))
        messageId := streamRes[0].Messages[0].ID
        
        //创建订单
        err = s.db.Transaction(func(tx *gorm.DB) error {
            //6.扣减库存
            //6.1 这里在高并发情况下会出现超卖的情况
            result := tx.Model(&models.SecKillVoucher{}).
                Where("stock > 0").
                Update("stock", gorm.Expr("stock - ?", 1))
            if result.Error != nil || result.RowsAffected == 0 {
                return fmt.Errorf("扣减库失败 %s", result.Error)
            }
            //7.创建订单
            
            order := &models.VoucherOrder{
                UserId:    uint64(userId),
                VoucherId: uint64(voucherId),
            }
            if err = tx.Create(order).Error; err != nil {
                return fmt.Errorf("创建订单失败: %s", err)
            }
            //确认消息
            s.rdb.XAck(ctx, "stream.orders", "g1", messageId)
            log.Println("创建订单成功")
            return nil
        })
    }
}

func (s CouponService) ListenQueue(ctx context.Context) {
    for {
        select {
        case order := <-s.queue:
            go func(order *models.VoucherOrder) {
                log.Printf("listen queue - uid: %d, couponId: %d", order.UserId, order.VoucherId)
                //6.扣减库存
                //6.1 这里在高并发情况下会出现超卖的情况
                err := s.db.Transaction(func(tx *gorm.DB) error {
                    result := tx.Model(&models.SecKillVoucher{}).
                        Where("voucher_id = ? and stock > 0", order.VoucherId).
                        Update("stock", gorm.Expr("stock - ?", 1))
                    if result.Error != nil || result.RowsAffected == 0 {
                        return fmt.Errorf("扣减库失败 %s uid:%d, couponId: %d", result.Error, order.UserId, order.VoucherId)
                    }
                    //7.创建订单
                    if err := tx.Create(order).Error; err != nil {
                        return fmt.Errorf("创建订单失败: %s uid:%d, couponId: %d", err, order.UserId, order.VoucherId)
                    }
                    return nil
                })
                if err != nil {
                    log.Println(err)
                } else {
                    log.Printf("下单成功: uid:%d, couponId: %d", order.UserId, order.VoucherId)
                }
            }(order)
        case <-ctx.Done():
            return
        }
    }
}

func (s CouponService) secKillCouponV3(ctx context.Context, couponId int, userId int) error {
    //1.执行lua脚本来判断用户和是否下过单
    result, err := luaAtomicOrder.Run(ctx, s.rdb, []string{}, couponId, userId).Result()
    if err != nil {
        return fmt.Errorf("下单失败:%s", err)
    }
    //2.判断是否为1库存不足
    flag, _ := result.(int64)
    switch flag {
    case 1:
        //3.库存不足
        return errors.New(strconv.Itoa(couponId) + "库存不足")
    case 2:
        //4.判断是否为2已下单
        return errors.New(strconv.Itoa(userId) + "用户已下过单")
    default:
        //5. 异步下单 已经将数据发送至异步队列所以这里直接返回nil即可，另一个线程去队列中读取 然后下单即可
        return nil
    }
}

func (s CouponService) secKillCouponV2(ctx context.Context, couponId int, userId int) (*models.VoucherOrder, error) {
    //1.执行lua脚本来判断用户和是否下过单
    result, err := luaAtomicOrder.Run(ctx, s.rdb, []string{}, couponId, userId).Result()
    if err != nil {
        return nil, fmt.Errorf("下单失败:%s", err)
    }
    //2.判断是否为1库存不足
    flag, _ := result.(int64)
    switch flag {
    case 1:
        //3.库存不足
        return nil, errors.New(strconv.Itoa(couponId) + "库存不足")
    case 2:
        //4.判断是否为2已下单
        return nil, errors.New(strconv.Itoa(userId) + "用户已下过单")
    default:
        //5.下单成功
        //5.1 生成订单id
        //5.2 创建订单
        order := &models.VoucherOrder{
            UserId:    uint64(userId),
            VoucherId: uint64(couponId),
        }
        //5.2 异步下单
        return order, err
    }
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
    var tryLock = &RedisLock{KeyPrefix: "biz_lock:", rdb: rdb}
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
