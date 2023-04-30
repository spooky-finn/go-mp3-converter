package api

import (
	"3205.team/go-mp3-converter/helpers"
	"3205.team/go-mp3-converter/redis"
	rdriver "github.com/redis/go-redis/v9"
)

type RedisAPI struct {
	client            *rdriver.Client
	taskHandlerParams *redis.TaskHandlerParams
}

func NewRedisAPI() *RedisAPI {
	return &RedisAPI{
		client: redis.Rdb,
		taskHandlerParams: &redis.TaskHandlerParams{
			QueueTimeout: 45,
		},
	}
}

func (r *RedisAPI) Init() {
	pooler := redis.NewPooler()
	taskHandler := redis.NewTaskHandler(r.taskHandlerParams)

	go func() {
		for task := range pooler.Ch {
			progress := helpers.NewProg()
			if err := taskHandler.Handle(task, progress); err != nil {
				panic(err)
			}
		}
	}()

}
