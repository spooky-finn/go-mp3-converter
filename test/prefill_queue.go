package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/pkg"

	"3205.team/go-mp3-converter/infra/redisscheduler"
)

var payload = []redisscheduler.Request{
	{
		OriginalURL: "https://www.youtube.com/watch?v=UBr3fOBvotc&list=RDUBr3fOBvotc&start_radio=1&ab_channel=Asoon",
		SourceURL:   "https://www.youtube.com/watch?v=UBr3fOBvotc&list=RDUBr3fOBvotc&start_radio=1&ab_channel=Asoon",
		PushedAt:    time.Now().Unix(),
	},
	// {
	// 	OriginalURL: "https://www.youtube.com/watch?v=dO3CaZb3I_s&ab_channel=LostSounds",
	// 	SourceURL:   "https://www.youtube.com/watch?v=dO3CaZb3I_s&ab_channel=LostSounds",
	// 	PushedAt:    time.Now().Unix(),
	// },
}

func main() {
	ctx := context.Background()

	for _, each := range payload {
		buf, err := json.Marshal(each)
		if err != nil {
			panic(err)
		}

		pkg.GetRedisClient().LPush(ctx, cfg.AppConfig.Rdb.QueueTable, string(buf))
	}

	fmt.Println("redis task queue fullfilled")
}
