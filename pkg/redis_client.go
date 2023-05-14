package pkg

import (
	"context"
	"fmt"

	"3205.team/go-mp3-converter/cfg"
	rdriver "github.com/redis/go-redis/v9"
)

var redisClient *rdriver.Client

func GetRedisClient() *rdriver.Client {
	if redisClient != nil {
		return redisClient
	}
	addr := fmt.Sprintf("%s:%d", cfg.AppConfig.Rdb.Host, cfg.AppConfig.Rdb.Port)
	c := rdriver.NewClient(
		&rdriver.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
			OnConnect: func(ctx context.Context, conn *rdriver.Conn) error {
				Logger.Printf("established connection with redis %v", addr)
				ping := conn.Ping(ctx)
				Logger.Printf("redis ping: %v \n", ping.Val())
				return nil
			},
		})

	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		Logger.Panicf("failed to connect to the redis: %v:", err)
	}

	redisClient = c
	return c
}
