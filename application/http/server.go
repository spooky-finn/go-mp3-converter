package http

import (
	"log"

	"3205.team/go-mp3-converter/domain/mp3converter"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type HTTPServer struct {
	addr         *string
	mp3converter *mp3converter.MP3Converter
}

func NewHTTPServer(addr *string, mp3Converter *mp3converter.MP3Converter) *HTTPServer {
	hs := &HTTPServer{
		addr:         addr,
		mp3converter: mp3Converter,
	}
	go hs.init()
	return hs
}

func (h *HTTPServer) init() {
	router := router.New()

	router.GET("/health", h.HandleHealth)
	router.POST("/api/mp3convert", h.HandleConvertToMP3)
	router.GET("/api/download/{params}", h.HandleDownload)

	server := &fasthttp.Server{
		Handler: func() fasthttp.RequestHandler {
			log.Printf("server ready to handle request at: %s", *h.addr)
			return router.Handler
		}(),
		IdleTimeout: 10000,
	}

	go func() {
		if err := server.ListenAndServe(*h.addr); err != nil {
			log.Fatalf("Error in ListenAndServe: %v", err)
		}
	}()
}
