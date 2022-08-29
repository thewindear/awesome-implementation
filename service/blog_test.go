package service

import (
    "implementation-scheme/models"
    "log"
    "testing"
)

func TestBlogLiked(t *testing.T) {
    service := &BlogService{db: db, rdb: rdb}
    err := service.LikeBlog(4, 1)
    if err != nil {
        t.Errorf("点赞失败:%s", err.Error())
    } else {
        t.Log("点赞成功")
    }
}

func TestBlogService_MessageBox(b *testing.T) {
    service := &BlogService{db: db, rdb: rdb}
    log.Println(service.MessageBox(ctx, 2, 1661796599, 2))
}

func TestBlogService_SaveBlog(t *testing.T) {
    service := &BlogService{db: db, rdb: rdb}
    blog := &models.Blog{
        ShopId:  9,
        UserId:  1,
        Title:   "hello world",
        Content: "今天是个好日子",
    }
    err := service.SaveBlog(ctx, blog)
    if err != nil {
        t.Error(err)
    } else {
        t.Log("发布成功")
    }
}

func TestBlogService_GetBlogById(t *testing.T) {
    service := &BlogService{db: db, rdb: rdb}
    blog, err := service.GetBlogById(ctx, 4, 7779412)
    if err != nil {
        t.Error(err)
    } else {
        t.Log(blog.IsLiked)
        t.Log(blog.LikeList)
    }
}

func TestBlogService_UserLikeBlog(t *testing.T) {
    service := &BlogService{db: db, rdb: rdb}
    err := service.UserLikeBlog(ctx, 7779410, 4)
    if err != nil {
        t.Errorf(err.Error())
    }
}
