package redisscheduler

import (
	"errors"

	"3205.team/go-mp3-converter/application/progress"
	"3205.team/go-mp3-converter/domain/mp3converter"
	"3205.team/go-mp3-converter/infra/cache"
	"3205.team/go-mp3-converter/pkg"
	"github.com/redis/go-redis/v9"
)

type RedisScheduler struct {
	taskhandler *Handler
	pooler      *Pooler
	cache       *cache.Cache
}

func NewRedisScheduler(redisclient *redis.Client, mp3converter *mp3converter.MP3Converter, cache *cache.Cache) *RedisScheduler {
	rs := &RedisScheduler{
		taskhandler: NewHandler(mp3converter),
		pooler:      NewPooler(redisclient),
		cache:       cache,
	}
	go rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	for request := range r.pooler.Ch {
		go r.handleIncomingRequest(request)
	}
}

func (r *RedisScheduler) handleIncomingRequest(request *Request) {
	prog := progress.New()
	_, err := r.taskhandler.Handle(request, prog)
	if err != nil {
		if errors.Is(err, ErrQueueTimeout) {
			pkg.Logger.Println("queue timeout elapsed")
		} else {
			panic(err)
		}
	}

	for range prog.Ch {
		println("progress in redis scheduler")
		continue
	}

}
