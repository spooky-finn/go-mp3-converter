package redisscheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/entity"
	rdriver "github.com/redis/go-redis/v9"
)

type Pusher struct {
	table  string
	client *rdriver.Client
}

func NewPusher(client *rdriver.Client) *Pusher {
	return &Pusher{
		table:  cfg.AppConfig.Rdb.TaskTable,
		client: client,
	}
}

func getKey(ID string) string {
	return fmt.Sprintf("%s:%s", cfg.AppConfig.Rdb.TaskTable, ID)
}

func (p *Pusher) PushTask(task *entity.Task, expiration time.Duration) {
	buf, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	p.client.Set(ctx, getKey(task.ID), buf, expiration)
}

// write a function to update task
func (p *Pusher) UpdateTask(task *entity.Task, expiration time.Duration) {
	buf, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	p.client.Set(ctx, getKey(task.ID), buf, expiration)
}
