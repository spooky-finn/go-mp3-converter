package redisscheduler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"3205.team/go-mp3-converter/cfg"
	rdriver "github.com/redis/go-redis/v9"
)

type Pooler struct {
	table  string
	client *rdriver.Client
	Ch     chan *IncomingTask
}

var ctx = context.Background()
var NullTask = &IncomingTask{}

func NewPooler(client *rdriver.Client) *Pooler {
	p := &Pooler{
		table:  cfg.AppConfig.Rdb.QueueTable,
		client: client,
		Ch:     make(chan *IncomingTask),
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
		request, err := p.pool()
		if err != nil {
			if errors.Is(err, &json.UnmarshalTypeError{}) {
				panic(err)
			}
			continue
		}
		p.Ch <- request
	}
}

func (p *Pooler) pool() (*IncomingTask, error) {
	var j = p.client.LPop(ctx, p.table)
	if err := j.Err(); err != nil {
		return NullTask, err
	}

	request := &IncomingTask{}
	err := json.Unmarshal([]byte(j.Val()), request)
	if err != nil {
		return NullTask, err
	}

	return request, nil
}

func (p *Pooler) Push(task *OutgoingTask, ttl time.Duration) error {
	b, err := json.Marshal(task)
	if err != nil {
		return err
	}

	key := cfg.AppConfig.Rdb.GetCachedTaskKey(task.TaskName)

	return p.client.Set(ctx, key, b, ttl).Err()
}
