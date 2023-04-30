package cfg

import (
	"os"
	"time"
)

type RedisConfig struct {
	Host         string
	Port         int
	TaskTable    string
	QueueTable   string
	PoolInterval time.Duration
}

type ConfigType struct {
	FfmpegbinPath  string
	FfprobebinPath string
	// directory to save converted files and currenty downloading files
	TempDir string
	// time to live for converted file
	// if file is not downloaded in this time, it will be deleted
	FileTTL time.Duration
	Rdb     RedisConfig
	// link microservice
	LinkURL string
}

var AppConfig = &ConfigType{
	FfmpegbinPath:  "./bin/ffmpeg",
	FfprobebinPath: "./bin/ffprobe",
	TempDir:        "./tmp",
	FileTTL:        time.Minute * 15,
	Rdb: RedisConfig{
		Host:         "localhost",
		Port:         6379,
		TaskTable:    "tasks",
		QueueTable:   "mp3Queue",
		PoolInterval: time.Millisecond * 100,
	},
	LinkURL: os.Getenv("LINK_URL"),
}
