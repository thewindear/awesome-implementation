package service

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "implementation-scheme/models"
    "log"
    "time"
)

type FollowService struct {
    rdb *redis.Client
    db  *gorm.DB
}

const (
    followListKey = "followList:%d"
)

// CommonFollow 共同关注的用户集合
func (f FollowService) CommonFollow(ctx context.Context, operationOf uint64, userId uint64) ([]string, error) {
    key1 := fmt.Sprintf(followListKey, operationOf)
    key2 := fmt.Sprintf(followListKey, userId)
    commons := f.rdb.ZInter(ctx, &redis.ZStore{
        Keys: []string{key2, key1},
    }).Val()
    return commons, nil
}

// Follow 关注
func (f FollowService) Follow(ctx context.Context, operationOf uint64, followUid uint64) error {
    key := fmt.Sprintf(followListKey, operationOf)
    // 1. 判断用户是否关注过followUId
    var err error
    exists := f.rdb.ZScore(ctx, key, fmt.Sprintf("%d", followUid)).Val()
    //1.1已经关注
    if exists > 0 {
        // 1.2 取消关注
        log.Printf("%d 取消关注 %d", operationOf, followUid)
        return f.unFollow(ctx, operationOf, followUid)
    } else {
        // 1.3 未关注 关注用户并添加数据库
        follow := &models.Follow{
            UserId:       operationOf,
            FollowUserId: followUid,
        }
        log.Printf("%d 关注 %d", operationOf, followUid)
        err = f.db.Create(follow).Error
        if err == nil {
            err = f.rdb.ZAdd(ctx, key, redis.Z{
                Member: followUid,
                Score:  float64(time.Now().Unix()),
            }).Err()
        }
        return err
    }
    
}

// UnFollow 取关
func (f FollowService) unFollow(ctx context.Context, operationOf uint64, unFollowUid uint64) error {
    err := f.db.Where("user_id = ? AND follow_user_id = ?", operationOf, unFollowUid).
        Delete(&models.Follow{}).Error
    if err != nil {
        return nil
    }
    key := fmt.Sprintf(followListKey, operationOf)
    return f.rdb.ZRem(ctx, key, unFollowUid).Err()
}
