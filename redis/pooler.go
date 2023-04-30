package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"3205.team/go-mp3-converter/cfg"
	redis "github.com/redis/go-redis/v9"
)

type Status string

const (
	StatusNew      Status = "new"
	StatusReady    Status = "ready"
	StatusError    Status = "error"
	StatusProgress Status = "in_progress"
)

type Pooler struct {
	TaskTable string
	Client    *redis.Client
	Ch        chan *Task
}

var ctx = context.Background()
var logger = log.New(os.Stdout, "redis: ", log.LstdFlags)

func NewPooler() *Pooler {
	p := &Pooler{
		TaskTable: cfg.AppConfig.Rdb.TaskTable,
		Client:    Rdb,
		Ch:        make(chan *Task),
	}

	// run getTask in goroutine periodically every 50ms
	// and send task to channel
	go func() {
		defer p.Client.Close()
		for {
			<-time.After(cfg.AppConfig.Rdb.PoolInterval)
			task, err := p.poolTask()

			if err != nil {
				logger.Println(err)
				continue
			}
			p.Ch <- task
		}
	}()

	return p
}

func (p *Pooler) poolTask() (*Task, error) {
	var j = p.Client.LPop(ctx, p.TaskTable)
	fmt.Printf("pooled task: %+v\n", j)
	if err := j.Err(); err != nil {
		return &Task{}, err
	}

	task := &Task{
		Status: StatusNew,
	}

	err := json.Unmarshal([]byte(j.Val()), &Task{})
	if err != nil {
		return &Task{}, err
	}

	return task, nil
}
