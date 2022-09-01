package service

import (
    "encoding/json"
    "fmt"
    "log"
    "testing"
)

func TestPreheatData(t *testing.T) {
    service := NewItemService()
    items, _ := service.FindAllItem()
    for _, item := range items {
        itemJson, _ := json.Marshal(item)
        cacheKey := fmt.Sprintf("item:%d", item.Id)
        rdb.Set(ctx, cacheKey, itemJson, 0)
    }
    stocks, _ := service.FindAllItemStock()
    for _, item := range stocks {
        itemJson, _ := json.Marshal(item)
        cacheKey := fmt.Sprintf("stock:%d", item.ItemId)
        rdb.Set(ctx, cacheKey, itemJson, 0)
    }
    
}

func TestItemService_GetById(t *testing.T) {
    service := NewItemService()
    log.Println(service.GetById(10002))
}
