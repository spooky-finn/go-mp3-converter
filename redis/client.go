package redis

import (
	"fmt"

	"3205.team/go-mp3-converter/cfg"
	"github.com/redis/go-redis/v9"
)

var Rdb = redis.NewClient(
	&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.AppConfig.Rdb.Host, cfg.AppConfig.Rdb.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
