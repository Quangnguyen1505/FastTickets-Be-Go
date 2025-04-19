package initialize

import (
	"context"
	"fmt"

	"github.com/ntquang/ecommerce/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Db,
		PoolSize: r.PoolSize,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Failed to connect to Redis, error::", zap.Error(err))
		return
	}

	global.Logger.Info("initialization Redis Successfully!")
	global.Redis = rdb
}
