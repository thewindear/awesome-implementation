package handler

import (
    "github.com/gofiber/fiber/v2"
    "implementation-scheme/service"
)

var itemShop *service.ItemService

func init() {
    itemShop = service.NewItemService()
}

func ShopItem(ctx *fiber.Ctx) error {
    id, _ := ctx.ParamsInt("id", 0)
    if id == 0 {
        return fiber.NewError(400, "商品id不能为空")
    }
    item, err := itemShop.GetById(uint64(id))
    if err != nil {
        return fiber.NewError(500, "服务异步")
    }
    if item == nil {
        return fiber.NewError(404, "资源不存在")
    }
    return ctx.JSON(item)
}

func ShopItemStock(ctx *fiber.Ctx) error {
    id, _ := ctx.ParamsInt("id", 0)
    if id == 0 {
        return fiber.NewError(400, "商品id不能为空")
    }
    item, err := itemShop.GetStockById(uint64(id))
    if err != nil {
        return fiber.NewError(500, "服务异步")
    }
    if item == nil {
        return fiber.NewError(404, "资源不存在")
    }
    return ctx.JSON(item)
}
