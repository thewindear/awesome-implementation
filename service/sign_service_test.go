package service

import "testing"

func TestSign(t *testing.T) {
    service := &SignService{rdb: rdb, db: db}
    t.Log(service.Sign(32131))
    t.Log(service.SignCount(32131))
}
