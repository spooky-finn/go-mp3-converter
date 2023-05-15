package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/pkg"

	"3205.team/go-mp3-converter/entity"

	rdriver "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	client       *rdriver.Client
	ttl          time.Duration
	disableCache bool
}

func NewCache() *Cache {
	return &Cache{
		client:       pkg.GetRedisClient(),
		ttl:          1 * time.Hour,
		disableCache: os.Getenv("DISABLE_CACHE") == "true",
	}
}

func getKey(ID string) string {
	return cfg.AppConfig.Rdb.TaskTable + "." + ID
}

func (c *Cache) getTaskByKey(key string) *entity.Task {
	buf, err := c.client.Get(ctx, key).Bytes()
	if err == rdriver.Nil {
		return nil
	} else if err != nil {
		panic(err)
	}

	task := &entity.Task{}
	if err := json.Unmarshal(buf, task); err != nil {
		panic(err)
	}

	task.WasCached = true
	return task
}

func (c *Cache) GetTask(ID string) *entity.Task {
	if c.disableCache {
		return nil
	}

	return c.getTaskByKey(getKey(ID))
}

func (c *Cache) SetTask(task *entity.Task) {
	buf, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	if err := c.client.Set(ctx, getKey(task.ID), buf, c.ttl).Err(); err != nil {
		panic(err)
	}
}
