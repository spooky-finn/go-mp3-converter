package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"3205.team/go-mp3-converter/application/redisscheduler"
	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain/mp3converter"
	"3205.team/go-mp3-converter/pkg"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var payload = []redisscheduler.IncomingTask{
	{
		TaskName:    "task1",
		OriginalURL: "https://www.youtube.com/watch?v=UBr3fOBvotc&list=RDUBr3fOBvotc&start_radio=1&ab_channel=Asoon",
		SourceURL:   "https://www.youtube.com/watch?v=UBr3fOBvotc&list=RDUBr3fOBvotc&start_radio=1&ab_channel=Asoon",
		Origin:      "example.com",
		PushedAt:    time.Now().Unix(),
	},
	// {
	// 	TaskName:    "task2",
	// 	OriginalURL: "https://www.youtube.com/watch?v=dO3CaZb3I_s&ab_channel=LostSounds",
	// 	SourceURL:   "https://www.youtube.com/watch?v=dO3CaZb3I_s&ab_channel=LostSounds",
	// 	Origin:      "example.com",
	// 	PushedAt:    time.Now().Unix(),
	// },
	{
		TaskName:    "task3",
		OriginalURL: "https://www.youtube.com/watch?v=69JaMvxmhCA&list=RDPWpYpTFLoK0&index=4&ab_channel=coldcarti-Topic",
		SourceURL:   "https://www.youtube.com/watch?v=69JaMvxmhCA&list=RDPWpYpTFLoK0&index=4&ab_channel=coldcarti-Topic",
		Origin:      "example.com",
		PushedAt:    time.Now().Unix(),
	},
}

func clearQueue(client *redis.Client) {
	ctx := context.Background()
	client.Del(ctx, cfg.AppConfig.Rdb.QueueTable)
}

func prefillQueue(client *redis.Client) {
	ctx := context.Background()

	for _, each := range payload {
		buf, err := json.Marshal(each)
		if err != nil {
			panic(err)
		}

		client.LPush(ctx, cfg.AppConfig.Rdb.QueueTable, string(buf))
	}
	fmt.Println("redis task queue fullfilled")
}

// func TestFullyPrefillQueue(t *testing.T) {
// 	client := pkg.GetRedisClient()
// 	ctx := context.Background()

// 	ClearQueue(client)
// 	PrefillQueue(client)

// 	itemsInRedis := client.LLen(ctx, cfg.AppConfig.Rdb.QueueTable)
// 	assert.Equal(t, int64(3), itemsInRedis.Val(), "redis task queue not fullfilled")
// }

func TestRedisScheduler(t *testing.T) {
	tempDir := cfg.AppConfig.TempDir
	os.Mkdir(tempDir, os.ModePerm)
	err := godotenv.Load(".env")
	if err != nil {
		pkg.Logger.Fatalf("Error loading .env file: %v", err)
	}

	client := pkg.GetRedisClient()
	ctx := context.Background()
	mp3converter := mp3converter.New()

	clearQueue(client)
	prefillQueue(client)

	temsInRedis := client.LLen(ctx, cfg.AppConfig.Rdb.QueueTable)
	assert.Equal(t, int64(len(payload)), temsInRedis.Val(), "redis task queue not empty")

	// create redis scheduler
	redisscheduler.NewRedisScheduler(client, mp3converter)

	// wait for all tasks to be processed
	time.Sleep(3 * time.Minute)
}
