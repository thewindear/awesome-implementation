package service

import (
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
