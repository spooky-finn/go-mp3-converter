package redisscheduler

import (
	"errors"

	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"

	rdriver "github.com/redis/go-redis/v9"
)

type RedisScheduler struct {
	client *rdriver.Client
}

func NewRedisScheduler() *RedisScheduler {
	rs := &RedisScheduler{
		client: pkg.GetRedisClient(),
	}
	rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	pooler := NewPooler(r.client)
	taskhandler := NewHandler()

	go func() {
		for request := range pooler.Ch {
			task := entity.NewTask(entity.NewTaskParams{
				OriginalURL: request.OriginalURL,
				// DownloadURL: request.DownloadURL,
			})

			go func() {
				for v := range task.Progress.Ch {
					pkg.Logger.Printf("progress: %v%%\n", v)
				}
			}()

			err := taskhandler.Handle(task)
			if err != nil {
				if errors.Is(err, ErrQueueTimeout) {
					pkg.Logger.Println("queue timeout elapsed")
				} else {
					panic(err)
				}
			}

			// key := fmt.Sprintf("%s:%s", cfg.AppConfig.Rdb.TaskTable, ce.Filename)
			// buf, err := json.Marshal(ce)
			// if err != nil {
			// 	panic(err)
			// }

			// r.client.Set(ctx, key, buf, 0)

		}
	}()

}
