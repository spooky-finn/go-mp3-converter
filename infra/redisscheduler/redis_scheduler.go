package redisscheduler

import (
	"errors"

	"3205.team/go-mp3-converter/pkg"
)

type RedisScheduler struct {
	taskhandler *Handler
	pooler      *Pooler
	pusher      *Pusher
}

func NewRedisScheduler() *RedisScheduler {
	redisclient := pkg.GetRedisClient()

	rs := &RedisScheduler{
		taskhandler: NewHandler(),
		pooler:      NewPooler(redisclient),
		pusher:      NewPusher(redisclient),
	}
	rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	for request := range r.pooler.Ch {
		println("new request")
		go r.handleIncomingRequest(request)
	}

}

func (r *RedisScheduler) handleIncomingRequest(request *Request) {
	task, err := r.taskhandler.Handle(request)
	if err != nil {
		if errors.Is(err, ErrQueueTimeout) {
			pkg.Logger.Println("queue timeout elapsed")
		} else {
			panic(err)
		}
	}

	<-task.Done
	r.pusher.Push(task, 0)
}
