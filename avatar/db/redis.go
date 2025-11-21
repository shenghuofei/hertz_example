package db

import (
	"avatar/config"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	rs     *redsync.Redsync
	once   sync.Once
)

// Init 初始化 Redis 客户端
func InitRedis() {
	once.Do(func() {
		addr := config.Cfg.GetString("redis.addr")
		client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     config.Cfg.GetString("redis.passwd"),
			DB:           config.Cfg.GetInt("redis.db"),
			MaxRetries:   config.Cfg.GetInt("redis.max_retries"),
			DialTimeout:  time.Duration(config.Cfg.GetInt("redis.dial_timeout_sec")) * time.Second,
			ReadTimeout:  time.Duration(config.Cfg.GetInt("redis.read_timeout_sec")) * time.Second,
			WriteTimeout: time.Duration(config.Cfg.GetInt("redis.write_timeout_sec")) * time.Second,
			PoolSize:     config.Cfg.GetInt("redis.pool_size"),
			MinIdleConns: config.Cfg.GetInt("redis.min_idle_conns"),
		})

		// 初始化分布式锁
		pool := goredis.NewPool(client)
		rs = redsync.New(pool)

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

// GetRedisClient 返回全局 Redis 客户端
func GetRedisClient() *redis.Client {
	if client == nil {
		panic("redis not initialized")
	}
	return client
}

// GetRedsync 返回分布式锁实例
func GetRedsync() *redsync.Redsync {
	if rs == nil {
		panic("redis pool not initialized")
	}
	return rs
}
