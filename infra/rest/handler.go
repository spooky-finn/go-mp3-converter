package rest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain/converter"
	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/infra/cache"
	"3205.team/go-mp3-converter/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

func HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("ok")
}

type HandleConvertToMP3Req struct {
	// a direct url to the video file
	DownloadURL string `validate:"required,http_url" json:"url"`
	// a url for video (e.g to th youtube video)
	OriginalURL string `validate:"required,http_url" json:"originalUrl"`
}

type HandleConvertToMP3Result struct {
	Status   string `json:"status"` // "ok" or "error"
	Filename string `json:"filename,omitempty"`
}

func HandleConvertToMP3(ctx *fasthttp.RequestCtx) {
	pkg.Logger.SetPrefix("rest: ")

	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	resultchan := make(chan HandleConvertToMP3Result)
	var indata HandleConvertToMP3Req

	err := json.Unmarshal(ctx.Request.Body(), &indata)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(indata); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	task := entity.NewTask(entity.NewTaskParams{
		OriginalURL: indata.OriginalURL,
		DownloadURL: indata.DownloadURL,
	})

	pkg.Logger.Printf("new request for convertation from %s \n", ctx.RemoteAddr())
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for loop := true; loop; {
			select {
			case result := <-resultchan:
				buf, err := json.Marshal(result)
				if err != nil {
					panic(err)
				}

				fmt.Fprintf(w, "data: %s\n\n", string(buf))
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				w.Flush()
				if err := ctx.Response.CloseBodyStream(); err != nil {
					panic(err)
				}
				loop = false
			case prog, ok := <-task.Progress.Ch:
				pkg.Logger.Printf("prog: %v", prog)
				if !ok {
					continue
				}

				buf, err := json.Marshal(prog)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, "data: %s\n\n", string(buf))
				w.Flush()
			}

		}
	}))

	go func() {
		cache := cache.NewCache()
		cachedTask := cache.Get(task.ID)
		if cachedTask != nil {
			resultchan <- HandleConvertToMP3Result{
				Filename: cachedTask.ID,
				Status:   "ok",
			}
			return
		}

		err = converter.Run(task)
		if err != nil {
			resultchan <- HandleConvertToMP3Result{
				Status: "error",
			}
			return
		}
		resultchan <- HandleConvertToMP3Result{
			Filename: task.ID,
			Status:   "ok",
		}
		cache.Set(task)
	}()
}

// WIP
func HandleDownload(ctx *fasthttp.RequestCtx) {
	pkg.Logger.SetPrefix("rest: ")
	params := strings.Split(ctx.UserValue("params").(string), ".")

	if len(params) != 3 {
		ctx.Error("invalid params", fasthttp.StatusBadRequest)
		return
	}

	filename, timestamp, hash := params[0], params[1], params[2]

	pkg.Logger.Printf("new request for download from %s with filename: %s, timestamp: %s, hash: %s \n", ctx.RemoteAddr(), filename, timestamp, hash)
	fasthttp.ServeFile(ctx, fmt.Sprintf("%s/%s.mp3", cfg.AppConfig.TempDir, filename))
}
