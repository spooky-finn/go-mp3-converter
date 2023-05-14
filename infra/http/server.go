package http

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func NewRESTServer(addr *string) {
	router := router.New()

	router.GET("/health", HandleHealth)
	router.POST("/api/mp3convert", HandleConvertToMP3)
	router.GET("/api/download/{params}", HandleDownload)

	server := &fasthttp.Server{
		Handler: func() fasthttp.RequestHandler {
			log.Printf("server ready to handle request at: %s", *addr)
			return router.Handler
		}(),
		IdleTimeout: 10000,
	}

	go func() {
		if err := server.ListenAndServe(*addr); err != nil {
			log.Fatalf("Error in ListenAndServe: %v", err)
		}
	}()
}
