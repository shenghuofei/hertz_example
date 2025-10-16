package db

import (
	"avatar/config"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Init 初始化 Redis 客户端
func InitRedis() {
	once.Do(func() {
		addr := config.Cfg.GetString("redis.addr")
		client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     config.Cfg.GetString("redis.passwd"),    // 密码（如无可留空）
			DB:           config.Cfg.GetInt("redis.db"),           // 默认DB
			MinIdleConns: config.Cfg.GetInt("redis.MinIdleConns"), // 保持的最小空闲连接数
		})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("[redis] connect failed: %v", err))
		}
		hlog.Infof("[redis] connected to %s", addr)
	})
}

// CloseRedis 优雅关闭 Redis 客户端
func CloseRedis() {
	if client != nil {
		if err := client.Close(); err != nil {
			hlog.Errorf("[redis] close error: %v", err)
		} else {
			hlog.Info("[redis] connection closed")
		}
	}
}

// Client 返回全局 Redis 客户端
func RedisClient() *redis.Client {
	if client == nil {
		panic("redis not initialized")
	}
	return client
}
