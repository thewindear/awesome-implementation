package service

import (
    "log"
    "testing"
)

func TestFollowService_Follow(f *testing.T) {
    service := &FollowService{db: db, rdb: rdb}
    log.Println(service.Follow(ctx, 1, 7))
}

func TestFollowService_CommonFollow(t *testing.T) {
    service := &FollowService{db: db, rdb: rdb}
    log.Println(service.CommonFollow(ctx, 1, 7))
}
