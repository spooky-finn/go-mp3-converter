package redisscheduler

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/application/progress"
	"3205.team/go-mp3-converter/domain/mp3converter"
	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"
	"github.com/redis/go-redis/v9"
)

type RedisScheduler struct {
	taskhandler *Handler
	pooler      *Pooler
}

func NewRedisScheduler(redisclient *redis.Client, mp3converter *mp3converter.MP3Converter) *RedisScheduler {
	rs := &RedisScheduler{
		taskhandler: NewHandler(mp3converter),
		pooler:      NewPooler(redisclient),
	}
	go rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	for request := range r.pooler.Ch {
		go r.handleIncomingRequest(request)
	}
}

func (r *RedisScheduler) handleIncomingRequest(request *IncomingTask) {
	prog := progress.New()
	task, err := r.taskhandler.Handle(request, prog)
	if err != nil {
		if errors.Is(err, ErrQueueTimeout) {
			pkg.Logger.Println("queue timeout elapsed")
		} else {
			panic(err)
		}
	}

	for range prog.Ch {
		println("progress in redis scheduler")
		outgoingTask := convertToOutgoingTask(task, *request)

		if err := r.pooler.Push(outgoingTask, 1*time.Hour); err != nil {
			panic(err)
		}
		continue
	}

}

func convertToOutgoingTask(task *entity.Task, icomingtask IncomingTask) *OutgoingTask {
	return nil
}
