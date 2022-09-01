package main

import (
    "flag"
    "github.com/gofiber/fiber/v2"
    "implementation-scheme/handler"
    "log"
    "net/http"
)

func main() {
    var port string
    flag.StringVar(&port, "port", "8000", "server port")
    flag.Parse()
    
    app := fiber.New(fiber.Config{
        ErrorHandler: func(ctx *fiber.Ctx, err error) error {
            if err == nil {
                return ctx.Next()
            }
            var code = http.StatusInternalServerError
            if fiberErr, ok := err.(*fiber.Error); ok {
                code = fiberErr.Code
                return ctx.Status(code).JSON(fiberErr)
            } else {
                return ctx.Status(code).JSON(err)
            }
        },
    })
    app.Get("/item/:id", handler.ShopItem)
    app.Get("/item/:id/stock", handler.ShopItemStock)
    
    log.Fatalln(app.Listen(":" + port))
    
}
