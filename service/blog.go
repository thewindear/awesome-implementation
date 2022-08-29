package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "implementation-scheme/models"
    "log"
    "time"
)

type BlogService struct {
    rdb *redis.Client
    db  *gorm.DB
}

const (
    likeKey = "blog:liked:%d"
)

func (b BlogService) UserLikeBlog(ctx context.Context, userId, blogId uint64) (err error) {
    // 1. 判断用户是否已点赞
    key := fmt.Sprintf(likeKey, blogId)
    exists := b.rdb.ZScore(ctx, key, fmt.Sprintf("%d", userId)).Val()
    // 1.1 不存在没点赞
    if exists == 0 {
        err = b.LikeBlog(blogId, 1)
        if err != nil {
            return
        }
        b.rdb.ZAddNX(ctx, key, redis.Z{
            Score:  float64(time.Now().Unix()),
            Member: userId,
        })
        log.Println("点赞成功")
    } else {
        // 1.2 已经点赞 点赞数-1 删除 set集合数据
        err = b.LikeBlog(blogId, -1)
        if err != nil {
            return
        }
        b.rdb.ZRem(ctx, key, fmt.Sprintf("%d", userId))
        log.Println("取消点赞")
    }
    return nil
}

func (b BlogService) GetBlogById(ctx context.Context, blogId uint64, userId uint64) (blog *models.Blog, err error) {
    err = b.db.Model(&blog).First(&blog, "id = ?", blogId).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("blog不存在")
        }
        return nil, err
    }
    key := fmt.Sprintf(likeKey, blogId)
    exists := b.rdb.ZScore(ctx, key, fmt.Sprintf("%d", userId)).Val()
    if exists != 0 {
        blog.IsLiked = 1
    }
    blog.LikeList = b.rdb.ZRevRange(ctx, key, 0, 5).Val()
    return blog, nil
}

// LikeBlog 点赞
func (b BlogService) LikeBlog(id uint64, likeNum int8) (err error) {
    var blog models.Blog
    err = b.db.Model(&blog).First(&blog, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("blog不存在")
        }
        return fmt.Errorf("服务器异步: %s", err)
    }
    err = b.db.Model(&blog).
        Where("id = ?", blog.Id).
        Update("liked", gorm.Expr("liked + ?", likeNum)).Error
    if err != nil {
        return fmt.Errorf("点赞失败:%s", err)
    }
    return err
}
