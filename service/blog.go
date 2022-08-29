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
    feedKey = "feed:%d"
)

func (b BlogService) MessageBox(ctx context.Context, userId int64, lastId int64, offset int64) (blogs []string, nextLastId int64, nextOffset int64) {
    key := fmt.Sprintf(feedKey, userId)
    // 1.查询收件箱
    // ZREVRANGEBYSCORE test:sortedset 2 0 withscores limit 1 3
    if lastId == 0 {
        lastId = 1661797999 //time.Now().Unix()
        offset = 0
    }
    blogsZ := b.rdb.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
        Min:    "0",
        Max:    fmt.Sprintf("%d", lastId),
        Offset: offset,
        Count:  2,
    }).Val()
    log.Println(offset)
    if len(blogsZ) == 0 {
        nextLastId = lastId
        nextOffset = offset
        return
    }
    var minTime = 0.0
    for _, blogId := range blogsZ {
        blogs = append(blogs, blogId.Member.(string))
        if minTime == blogId.Score {
            nextOffset += 1
        } else {
            minTime = blogId.Score
            nextOffset = 1
        }
    }
    nextLastId = int64(minTime)
    return
}

func (b BlogService) SaveBlog(ctx context.Context, blog *models.Blog) error {
    err := b.db.WithContext(ctx).Model(blog).Create(blog).Error
    if err != err {
        return err
    }
    //1.查询笔记作者的所有粉丝
    var fansList []*models.Follow
    err = b.db.Model(&models.Follow{}).
        Select("user_id").
        Where("follow_user_id = ?", blog.UserId).Find(&fansList).Error
    if err != nil {
        return err
    }
    //2.推送笔记id给所有粉丝
    for _, fans := range fansList {
        fansFeedKey := fmt.Sprintf(feedKey, fans.UserId)
        b.rdb.ZAdd(ctx, fansFeedKey, redis.Z{
            Score:  float64(time.Now().Unix()),
            Member: blog.Id,
        })
    }
    return nil
}

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
