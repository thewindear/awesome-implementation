package service

import (
    "context"
    "database/sql"
    "github.com/go-redis/redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
    "sync"
    "time"
)

var rdb *redis.Client
var ctx = context.Background()
var db *gorm.DB
var sqlDB *sql.DB

const dsn = "root:Kb7DPGVY98Dv64S97M73gW7GKZjCusje@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
    rdb = redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        ReadTimeout:  time.Second * 10,
        WriteTimeout: time.Second * 10,
    })
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    sqlDB, _ = db.DB()
    sqlDB.SetMaxIdleConns(100)
    sqlDB.SetMaxOpenConns(500)
    sqlDB.SetConnMaxLifetime(time.Second * 10)
    sqlDB.SetConnMaxIdleTime(time.Minute * 5)
    if err != nil {
        log.Fatalln(err)
    }
    db = db.Debug()
}

func CheckPing() {
    _ = sqlDB.Ping()
}

func ConcurrenceFn(max int, fn func() error) {
    wg := sync.WaitGroup{}
    wg.Add(max)
    for i := 0; i < max; i++ {
        go func() {
            defer wg.Done()
            err := fn()
            if err != nil {
                log.Printf("下单失败: %s", err)
            } else {
                log.Printf("下单成功")
            }
        }()
    }
    wg.Wait()
    log.Println("concurrence done")
}
