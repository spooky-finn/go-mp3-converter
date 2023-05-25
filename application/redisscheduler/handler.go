package redisscheduler

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/application/progress"
	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"

	"3205.team/go-mp3-converter/domain/mp3converter"

	link "3205.team/go-mp3-converter/infra"
)

var (
	ErrTaskHandler  = errors.New("task handler error")
	ErrQueueTimeout = errors.Join(ErrTaskHandler, errors.New("queue timeout elapsed"))
)

// incoming task
type IncomingTask struct {
	TaskName    string `json:"taskName"`
	OriginalURL string `json:"originalUrl"`
	// just a url to video on youtube, instgram, etc. Its not a link to download file
	SourceURL string `json:"url"`
	PushedAt  int64  `json:"pushedAt"`
	Origin    string `json:"origin"`
}

type Handler struct {
	QueueTimeout int64
	mp3Converter *mp3converter.MP3Converter
}

type HadlerResult struct {
	Task         *entity.Task
	Err          error
	LinkResponse *link.LinkResultData
}

func NewHandler(mp3Converter *mp3converter.MP3Converter) *Handler {
	return &Handler{
		QueueTimeout: 3600,
		mp3Converter: mp3Converter,
	}
}

// send request to the link microservice to get a link to download the video
// and start convertation
func (th *Handler) Handle(r *IncomingTask, prog *progress.Progress) *HadlerResult {
	// check that from the moment of adding the task to the queue, the time has not expired
	if time.Now().Unix()-r.PushedAt > th.QueueTimeout {
		return &HadlerResult{
			Err: ErrQueueTimeout,
		}
	}

	linkResponse, err := link.Fetch(r.OriginalURL)
	if err != nil {
		pkg.Logger.Println("fetch error: ", err)
		return &HadlerResult{
			Err: ErrQueueTimeout,
		}
	}

	task := th.mp3Converter.StartConvertation(entity.NewTaskParams{
		OriginalURL: r.OriginalURL,
		DownloadURL: linkResponse.DownloadURL,
		TaskName:    r.TaskName,
	}, prog)

	return &HadlerResult{
		Task:         task,
		LinkResponse: linkResponse,
	}
}
