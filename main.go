package main

import (
	"flag"
	"log"
	"os"

	"3205.team/go-mp3-converter/handlers"
	"3205.team/go-mp3-converter/internal"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	os.Mkdir(internal.TEMP_DIR, os.ModePerm)

	flag.Parse()

	// h := requestHandler
	// if *compress {
	// 	h = fasthttp.CompressHandler(h)
	// }
	router := router.New()

	router.GET("/health", handlers.HandleHealth)

	router.POST("/api/mp3convert", handlers.HandleConvertToMP3)

	router.GET("/api/download/{params}", handlers.HandleDownload)

	server := &fasthttp.Server{
		Handler: func() fasthttp.RequestHandler {
			log.Printf("server ready to handle request at: %s", *addr)
			return router.Handler
		}(),
		IdleTimeout: 10000,
	}

	if err := server.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}
