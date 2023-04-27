package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"3205.team/go-mp3-converter/helpers"
	"3205.team/go-mp3-converter/internal"
	"github.com/valyala/fasthttp"
)

func HandleHealth(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("ok")
}

type HandleConvertToMP3Req struct {
	Url string `json:"url"`
}

type HandleConvertToMP3Result struct {
	Status   string `json:"status"` // "ok" or "error"
	Filename string `json:"filename,omitempty"`
}

func HandleConvertToMP3(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	resultchan := make(chan HandleConvertToMP3Result)

	progress := helpers.NewProg()
	var indata HandleConvertToMP3Req

	err := json.Unmarshal(ctx.Request.Body(), &indata)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
	}

	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for loop := true; loop; {
			select {
			case result := <-resultchan:
				json, err := json.Marshal(result)
				if err != nil {
					panic(err)
				}

				fmt.Fprintf(w, "data: %s\n\n", string(json))
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				w.Flush()
				if err := ctx.Response.CloseBodyStream(); err != nil {
					panic(err)
				}
				loop = false
			case prog, ok := <-progress.Ch:
				if !ok {
					continue
				}

				json, err := json.Marshal(prog)
				if err != nil {
					panic(err)
				}

				fmt.Fprintf(w, "data: %s\n\n", string(json))
				w.Flush()
			}
		}
	}))

	go func() {
		filename, err := internal.NewConvertEntity(indata.Url, progress).ConvertToMP3()
		if err != nil {
			resultchan <- HandleConvertToMP3Result{
				Status: "error",
			}
			return
		}

		resultchan <- HandleConvertToMP3Result{
			Filename: filename,
			Status:   "ok",
		}
	}()
}

func HandleDownload(ctx *fasthttp.RequestCtx) {
	params := strings.Split(ctx.UserValue("params").(string), ".")

	if len(params) != 3 {
		ctx.Error("invalid params", fasthttp.StatusBadRequest)
		return
	}

	filename, timestamp, hash := params[0], params[1], params[2]

	fmt.Fprintf(ctx, "task: %s, timestamp: %s, hash: %s", filename, timestamp, hash)

	fasthttp.ServeFile(ctx, fmt.Sprintf("%s/%s.mp3", internal.TEMP_DIR, filename))
}
