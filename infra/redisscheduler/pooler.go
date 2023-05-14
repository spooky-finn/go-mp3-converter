package redisscheduler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"3205.team/go-mp3-converter/cfg"
	rdriver "github.com/redis/go-redis/v9"
)

type Status string

const (
	StatusNew      Status = "new"
	StatusReady    Status = "ready"
	StatusError    Status = "error"
	StatusProgress Status = "in_progress"
)

type Pooler struct {
	table  string
	client *rdriver.Client
	Ch     chan *Request
}

var ctx = context.Background()
var NullTask = &Request{}

func NewPooler(client *rdriver.Client) *Pooler {
	p := &Pooler{
		table:  cfg.AppConfig.Rdb.QueueTable,
		client: client,
		Ch:     make(chan *Request),
	}

	go p.periodicallPuller()
	return p
}

// run getTask in goroutine periodically every "PoolInterval"ms
// and send task to channel
func (p *Pooler) periodicallPuller() {
	defer p.client.Close()
	for {
		<-time.After(cfg.AppConfig.Rdb.PoolInterval)
		task, err := p.pool()
		if err != nil {
			if errors.Is(err, &json.UnmarshalTypeError{}) {
				panic(err)
			}
			continue
		}
		p.Ch <- task
	}
}

func (p *Pooler) pool() (*Request, error) {
	var j = p.client.LPop(ctx, p.table)
	if err := j.Err(); err != nil {
		return NullTask, err
	}

	task := &Request{}
	err := json.Unmarshal([]byte(j.Val()), task)
	if err != nil {
		return NullTask, err
	}

	task.Status = StatusNew
	return task, nil
}
