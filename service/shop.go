package service

import (
    "fmt"
    "github.com/go-redis/redis/v9"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "implementation-scheme/models"
    "strconv"
)

type ShopService struct {
    rdb *redis.Client
    db  *gorm.DB
}

const (
    shopTypeKey = "shop:geo:%d"
)

// QueryShopByType 通过经纬度和type查询附近的商家
func (s ShopService) QueryShopByType(
    typeId uint64, current int, x float64, y float64) []*models.Shop {
    //1.判断是否需要根据坐标查询
    if x == 0 || y == 0 {
        var shops []*models.Shop
        s.db.Model(&models.Shop{}).Find(&shops, "type_id = ?", typeId)
        return shops
    }
    //2.计算分页参数
    from := (current - 1) * 5
    end := current * 5
    //3.查询redis 按照距离排序 分页 shopId ,distance
    shops := s.rdb.GeoSearchLocation(
        ctx,
        fmt.Sprintf(shopTypeKey, typeId),
        &redis.GeoSearchLocationQuery{
            GeoSearchQuery: redis.GeoSearchQuery{
                Longitude:  x,
                Latitude:   y,
                Radius:     5,
                RadiusUnit: "km",
                Sort:       "asc",
                Count:      end,
            },
            WithDist: true,
        },
    ).Val()
    //3.1截取 从from到end
    if len(shops) > 0 {
        shops = shops[from:end]
    }
    // 每个店铺距离
    var shopDistances = make(map[string]float64)
    var shopIds []int
    //4.解析出id
    for _, shop := range shops {
        shopDistances[shop.Name] = shop.Dist
        shopId, _ := strconv.Atoi(shop.Name)
        shopIds = append(shopIds, shopId)
    }
    var shopList []*models.Shop
    s.db.Model(&models.Shop{}).
        Clauses(clause.OrderBy{
            Expression: clause.Expr{
                SQL:                "FIELD(id, ?)",
                Vars:               []interface{}{shopIds},
                WithoutParentheses: true,
            },
        }).Where("id in (?)", shopIds).
        Find(&shopList)
    
    for _, shop := range shopList {
        shop.ShopDist = shopDistances[strconv.Itoa(int(shop.Id))]
    }
    //5.查询店铺
    return shopList
}
