package redisscheduler

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"

	"3205.team/go-mp3-converter/domain/mp3converter"

	"3205.team/go-mp3-converter/infra/cache"
	link "3205.team/go-mp3-converter/infra/microservice"
)

var (
	ErrTaskHandler  = errors.New("task handler error")
	ErrQueueTimeout = errors.Join(ErrTaskHandler, errors.New("queue timeout elapsed"))
)

type OutcomeTask struct {
	Progress    int    `json:"progress"`
	OriginalURL string `json:"originalUrl"`
	Status      Status `json:"status"`
	Origin      string `json:"origin"`
	PushedAt    int64  `json:"pushedAt"`
	Instance    string `json:"instance"`
	Meta        struct {
		Title    string `json:"title"`
		Source   string `json:"source"`
		Duration string `json:"duration"`
		Tags     string `json:"tags"`
	} `json:"meta"`
	Thumb    string  `json:"thumb"`
	FileSize int     `json:"fileSize"`
	StartAt  float64 `json:"startAt"`
	Duration int     `json:"duration"`
	Filename string  `json:"filename"`
	StopAt   float64 `json:"stopAt"`
}

// incoming task
type Request struct {
	OriginalURL string `json:"originalUrl"`
	// just a url to video on youtube, instgram, etc. Its not a link to download file
	SourceURL string `json:"url"`
	PushedAt  int64  `json:"pushedAt"`
	Status    Status `json:",omitempty"`
}

type Handler struct {
	QueueTimeout int64
	mp3useCase   mp3converter.MP3ConverterUseCase
}

func NewHandler() *Handler {
	cache := cache.NewCache()
	return &Handler{
		QueueTimeout: 3600,
		mp3useCase:   *mp3converter.New(cache),
	}
}

// send request to the link microservice to get a link to download the video
// and start convertation
func (th *Handler) Handle(r *Request) (*entity.Task, error) {
	// check that from the moment of adding the task to the queue, the time has not expired
	if time.Now().Unix()-r.PushedAt > th.QueueTimeout {
		return nil, ErrQueueTimeout
	}

	linkResp, err := link.Fetch(r.OriginalURL)
	if err != nil {
		pkg.Logger.Println("fetch error: ", err)
		return nil, err
	}

	task := th.mp3useCase.StartConvertation(mp3converter.ConverterParams{
		OriginalURL: r.OriginalURL,
	})

	task.DownloadURL = linkResp.DownloadURL
	task.Thumb = linkResp.Thumb

	return task, nil
}
