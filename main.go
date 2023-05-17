package main

import (
	"flag"
	"os"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain/mp3converter"
	"3205.team/go-mp3-converter/infra/cache"
	"3205.team/go-mp3-converter/pkg"

	tempcleaner "3205.team/go-mp3-converter/application"
	"3205.team/go-mp3-converter/application/http"
	"3205.team/go-mp3-converter/application/redisscheduler"

	"github.com/joho/godotenv"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	tempDir := cfg.AppConfig.TempDir
	os.Mkdir(tempDir, os.ModePerm)
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		pkg.Logger.Fatalf("Error loading .env file: %v", err)
	}

	checkBinaryExists()

	// creating application layer services
	redisClient := pkg.GetRedisClient()
	ttl := cfg.AppConfig.Rdb.TTL

	cache := cache.New(redisClient, ttl)
	tempcleaner.New(tempDir, ttl, 1*time.Minute)

	// creating domain use cases
	mp3converter := mp3converter.New(cache)

	// crateing controllers
	redisscheduler.NewRedisScheduler(redisClient, mp3converter, cache)
	http.NewHTTPServer(addr, mp3converter)

	// run forever
	select {}
}

func checkBinaryExists() {
	// check if ffmpeg and ffprobe binaries present in directory
	if _, err := os.Stat(cfg.AppConfig.FfmpegbinPath); os.IsNotExist(err) {
		panic("ffmpeg binary not found")
	}
	if _, err := os.Stat(cfg.AppConfig.FfprobebinPath); os.IsNotExist(err) {
		panic("ffprobe binary not found")
	}
}
