package redis

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/helpers"
	"3205.team/go-mp3-converter/internal"
	"3205.team/go-mp3-converter/internal/services"
)

var (
	ErrTaskHandler  = errors.New("task handler error")
	ErrQueueTimeout = errors.Join(ErrTaskHandler, errors.New("queue timeout elapsed"))
)

type Task struct {
	TaskName    string    `json:"taskName"`
	OriginalURL string    `json:"originalUrl"`
	SourceURL   string    `json:"url"`
	PushedAt    time.Time `json:"pushedAt"`
	Status      Status
}

type TaskHandlerParams struct {
	QueueTimeout float64
}

type TaskHandler struct {
	params *TaskHandlerParams
}

func NewTaskHandler(params *TaskHandlerParams) *TaskHandler {
	return &TaskHandler{
		params: params,
	}
}

func (th *TaskHandler) Handle(task *Task, prog *helpers.Progress) error {
	if time.Since(task.PushedAt).Seconds() > th.params.QueueTimeout {
		return ErrQueueTimeout
	}

	linkResp, err := services.LinkRequest(task.OriginalURL)
	if err != nil {
		return err
	}

	if _, err := internal.NewConvertEntity(linkResp.DownloadURL, prog).ConvertToMP3(); err != nil {
		return err
	}

	println("task handled")
	return nil
}
