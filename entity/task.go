package entity

import (
	"time"

	"3205.team/go-mp3-converter/pkg"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusReady      Status = "ready"
	StatusError      Status = "error"
	StatusInProgress Status = "in_progress"
)

// a main entity that represents a task for convertation to mp3
type Task struct {
	ID string `json:"filename"`
	// a url to the video (e.g. to the youtube video)
	OriginalURL string `json:"originalUrl"`
	// a url issued by Link microservice to download the video
	DownloadURL string  `json:"downloadUrl"`
	Status      Status  `json:"status,omitempty"`
	Duration    float64 `json:"duration"`
	Thumb       string  `json:"thumb"`

	PushedAt int64 `json:"pushedAt"`
	StartAt  int64 `json:"startAt"`
	StopAt   int64 `json:"stopAt"`
	// if error durtting convertation occurs, it will be stored here
	err error `json:"-"`

	// a channel that will be closed when convertation is done
	WasCached bool `json:"-"`
}

type NewTaskParams struct {
	OriginalURL string
	DownloadURL string
	Thumb       string
}

func NewTask(params NewTaskParams) *Task {
	// if params.OriginalURL == "" || params.DownloadURL == "" {
	// 	panic("empty params")
	// }

	ID := pkg.Hash(params.OriginalURL)
	return &Task{
		ID:          ID,
		OriginalURL: params.OriginalURL,
		Status:      StatusNew,
		DownloadURL: params.DownloadURL,
		Thumb:       params.Thumb,
	}
}

// make task sutisfy error interface
func (t *Task) Error() string {
	return t.err.Error()
}

func (t *Task) IsError() bool {
	return t.err != nil || t.Status == StatusError
}

// methow to call when convertation is done successfully or with error (in case of error, task will be marked as errored)
func (t *Task) Teardown(err error) {
	if err != nil {
		pkg.Logger.Printf("task %s: failed with error: %s", t.ID, err)
		t.err = err
		t.Status = StatusError
	} else {
		t.Status = StatusReady
	}

	t.StopAt = time.Now().Unix()
}
