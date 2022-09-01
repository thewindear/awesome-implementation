package service

import (
    "fmt"
    "github.com/go-redis/redis/v9"
    "github.com/patrickmn/go-cache"
    "gorm.io/gorm"
    "implementation-scheme/models"
    "time"
)

type ItemService struct {
    rdb   *redis.Client
    db    *gorm.DB
    cache *cache.Cache
}

func NewItemService() *ItemService {
    return &ItemService{rdb: rdb, db: db, cache: cache.New(5*time.Minute, 10*time.Minute)}
}

func (i ItemService) FindAllItem() ([]*models.Item, error) {
    var items []*models.Item
    err := i.db.Model(&models.Item{}).Where("status != 3").Find(&items).Error
    return items, err
}

func (i ItemService) FindAllItemStock() ([]*models.ItemStock, error) {
    var itemStocks []*models.ItemStock
    err := i.db.Model(&models.ItemStock{}).Find(&itemStocks).Error
    return itemStocks, err
}

func (i ItemService) GetById(id uint64) (*models.Item, error) {
    var item models.Item
    cacheKey := fmt.Sprintf("item:%d", id)
    data, exists := i.cache.Get(cacheKey)
    if !exists {
        //查询数据库
        err := i.db.Model(&item).Where("id = ? AND status != 3", id).First(&item).Error
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return nil, nil
            }
            return nil, err
        }
        i.cache.Set(cacheKey, &item, cache.NoExpiration)
        return &item, nil
    }
    return data.(*models.Item), nil
}

func (i ItemService) GetStockById(id uint64) (*models.ItemStock, error) {
    var stock models.ItemStock
    err := i.db.Model(&stock).Where("item_id = ?", id).First(&stock).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &stock, nil
}
