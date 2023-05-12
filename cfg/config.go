package cfg

import (
	"os"
	"time"
)

type RedisConfig struct {
	Host string
	Port int
	// table with converted files
	TaskTable    string
	QueueTable   string
	PoolInterval time.Duration
}

type ConfigType struct {
	FfmpegbinPath  string
	FfprobebinPath string
	// directory to save converted files and currenty downloading files
	TempDir          string
	ProgressInterval time.Duration
	// time to live for converted file
	// if file is not downloaded in this time, it will be deleted
	FileTTL time.Duration
	Rdb     RedisConfig
	// link microservice
	LinkURL string
}

var AppConfig = &ConfigType{
	FfmpegbinPath:    "./bin/ffmpeg",
	FfprobebinPath:   "./bin/ffprobe",
	TempDir:          "./tmp",
	ProgressInterval: time.Second * 1,
	FileTTL:          time.Minute * 15,
	Rdb: RedisConfig{
		Host:         "178.63.85.247",
		Port:         26233,
		TaskTable:    "ffmpeg:mp3:task",
		QueueTable:   "ffmpeg:mp3:queue",
		PoolInterval: time.Millisecond * 100,
	},
	LinkURL: os.Getenv("LINK_URL"),
}
