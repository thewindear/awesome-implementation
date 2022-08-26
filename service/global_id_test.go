package service

import (
	"testing"
)

func TestRedisIdWorker(t *testing.T) {
	global := NewGlobalId(32)
	id, err := global.CreateId("order")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("当前生成的id: %d", id)
	}

}
