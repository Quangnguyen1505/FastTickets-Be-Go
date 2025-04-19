package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ntquang/ecommerce/global"
	"github.com/redis/go-redis/v9"
)

func GetCache(ctx context.Context, key string, obj interface{}) error {
	data, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	} else if err == redis.Nil {
		return fmt.Errorf("key %s not found", key)
	}

	// convert json to object
	if err = json.Unmarshal([]byte(data), &obj); err != nil {
		return fmt.Errorf("failed convert json to object")
	}

	return nil
}

func GetHashCache(ctx context.Context, keyUser string, obj interface{}) error {
	JwtTokenRedis, err := global.Redis.HGetAll(ctx, keyUser).Result()
	if err != nil {
		return err
	} else if err == redis.Nil {
		return fmt.Errorf("key %s not found", keyUser)
	}

	// convert json to object
	jsonData, err := json.Marshal(JwtTokenRedis)
	if err != nil {
		return fmt.Errorf("failed to marshal map to JSON: %v", err)
	}

	if err = json.Unmarshal(jsonData, &obj); err != nil {
		return fmt.Errorf("failed to convert JSON to object: %v", err)
	}
	return nil
}
