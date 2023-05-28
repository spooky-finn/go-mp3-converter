package cfg

import (
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	Host string
	Port int
	// table with converted files
	TaskTable    string
	QueueTable   string
	PoolInterval time.Duration
	TTL          time.Duration
}

func (rc *RedisConfig) GetCachedTaskKey(ID string) string {
	return rc.TaskTable + "." + ID
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

func CrossPlatformBinary(path string) string {
	if os.PathSeparator == '\\' {
		return path + ".exe"
	}
	return path
}

var port, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))

var AppConfig = &ConfigType{
	FfmpegbinPath:    CrossPlatformBinary("./bin/ffmpeg"),
	FfprobebinPath:   CrossPlatformBinary("./bin/ffprobe"),
	TempDir:          "./tmp",
	ProgressInterval: time.Second * 1,
	FileTTL:          time.Minute * 15,
	Rdb: RedisConfig{
		Host:         os.Getenv("REDIS_HOST"),
		Port:         port,
		TaskTable:    "ffmpeg:mp3:task",
		QueueTable:   "ffmpeg:mp3:queue",
		PoolInterval: time.Millisecond * 100,
		TTL:          time.Hour * 1,
	},
	LinkURL: os.Getenv("LINK_URL"),
}
