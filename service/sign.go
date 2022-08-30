package service

import (
    "errors"
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "log"
    "time"
)

type SignService struct {
    rdb *redis.Client
    db  *gorm.DB
}

const (
    signKey = "sign:%d:%d:%d"
)

func (s SignService) SignCount(uid uint) int {
    now := time.Now()
    //1.1获取年
    year := now.Year()
    //1.2获取月
    month := int(now.Month())
    //1.3 获取当前日
    dayOfMonth := now.Day()
    key := fmt.Sprintf(signKey, uid, year, month)
    //2.获取本月截止今天为止所有签到记录，返回是一个十进制数
    tmp := s.rdb.BitField(ctx, key, "GET", fmt.Sprintf("u%d", dayOfMonth), "0").Val()
    if len(tmp) == 0 {
        return 0
    }
    log.Println(tmp[0])
    //2.1循环遍历
    bitNum := tmp[0]
    var count = 0
    for {
        if (bitNum & 1) == 0 {
            //4.如果为0说明 没签到 结束
            break
        } else {
            //4.1如果不为0说明已签到，计数器+1
            count += 1
        }
        //3.让这个数字与1做与运算，得到数字最后一个bit位,将最后一位抛弃判断下一位
        bitNum = bitNum >> 1
    }
    //4.2把数字右移1位，抛弃最后位
    return count
}

func (s SignService) Sign(uid uint) error {
    now := time.Now()
    //1.1获取年
    year := now.Year()
    //1.2获取月
    month := int(now.Month())
    //1.3 获取当前日
    dayOfMonth := now.Day()
    //1.计算签到key
    key := fmt.Sprintf(signKey, uid, year, month)
    dayOfSignOffset := int64(dayOfMonth - 1)
    if s.rdb.GetBit(ctx, key, dayOfSignOffset).Val() == 0 {
        s.rdb.SetBit(ctx, key, dayOfSignOffset, 1)
        return nil
    } else {
        return errors.New("今天您已经签到")
    }
    
}
