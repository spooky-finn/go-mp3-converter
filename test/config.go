package test

import (
	"os"
	"time"

	"3205.team/go-mp3-converter/cfg"
)

var TestConfig = &cfg.ConfigType{
	FfmpegbinPath:    cfg.CrossPlatformBinary("../bin/ffmpeg"),
	FfprobebinPath:   cfg.CrossPlatformBinary("../bin/ffprobe"),
	TempDir:          "./tmp",
	ProgressInterval: time.Second * 1,
	FileTTL:          time.Minute * 15,
	Rdb: cfg.RedisConfig{
		Host:         "178.63.85.247",
		Port:         26233,
		TaskTable:    "ffmpeg:mp3:task",
		QueueTable:   "ffmpeg:mp3:queue",
		PoolInterval: time.Millisecond * 100,
		TTL:          time.Hour * 1,
	},
	LinkURL: os.Getenv("LINK_URL"),
}
