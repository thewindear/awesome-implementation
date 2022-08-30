package service

import (
    "fmt"
    "testing"
)

//  used_memory:2264224
//  used_memory_human:2.16M
func TestHyperLog(t *testing.T) {
    var values []string
    var j = 0
    for i := 0; i < 1000000; i++ {
        j = i % 1000
        if j == 999 {
            rdb.PFAdd(ctx, "hl2", values)
            values = []string{}
        } else {
            values = append(values, fmt.Sprintf("user_%d", i))
        }
    }
    t.Log(rdb.PFCount(ctx, "hl2").Val())
}
