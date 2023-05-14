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
	go rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	for request := range r.pooler.Ch {
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

	for {
		select {
		case <-task.Done:
			r.pusher.PushTask(task, 0)
			return
		case <-task.Progress.Ch:
			pkg.Logger.Println("task progress updated")
			continue
		}
	}

}
