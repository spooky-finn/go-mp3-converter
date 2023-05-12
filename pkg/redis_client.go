package pkg

import (
	"fmt"

	"3205.team/go-mp3-converter/cfg"
	rdriver "github.com/redis/go-redis/v9"
)

var redisClient *rdriver.Client

func GetRedisClient() *rdriver.Client {
	if redisClient != nil {
		return redisClient
	}
	return rdriver.NewClient(
		&rdriver.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.AppConfig.Rdb.Host, cfg.AppConfig.Rdb.Port),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
}
