package main

import (
	"flag"
	"os"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/pkg"

	"3205.team/go-mp3-converter/infra/http"
	"3205.team/go-mp3-converter/infra/redisscheduler"

	"github.com/joho/godotenv"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	os.Mkdir(cfg.AppConfig.TempDir, os.ModePerm)
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		pkg.Logger.Fatalf("Error loading .env file: %v", err)
	}

	checkBinaryExists()

	// init infrastructures
	redisscheduler.NewRedisScheduler()
	http.NewRESTServer(addr)

	// h := requestHandler
	// if *compress {
	// 	h = fasthttp.CompressHandler(h)
	// }
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
