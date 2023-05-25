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

type TaskTimePoints struct {
	StartAt int64
	StopAt  int64
}

// a main entity that represents a task for convertation to mp3
type Task struct {
	ID string
	// a url to the video (e.g. to the youtube video)
	OriginalURL string
	// a url issued by Link microservice to download the video
	DownloadURL string
	Status      Status
	Duration    float64

	FileSize int64
	// if error durtting convertation occurs, it will be stored here
	Err error

	TaskTimePoints
}

type NewTaskParams struct {
	TaskName    string
	OriginalURL string
	DownloadURL string
}

func NewTask(params NewTaskParams) *Task {
	// ID := pkg.Hash(params.OriginalURL)

	return &Task{
		ID:          params.TaskName,
		OriginalURL: params.OriginalURL,
		Status:      StatusNew,
		DownloadURL: params.DownloadURL,
	}
}

// make task sutisfy error interface
func (t *Task) Error() string {
	return t.Err.Error()
}

func (t *Task) IsError() bool {
	return t.Err != nil || t.Status == StatusError
}

// methow to call when convertation is done successfully or with error (in case of error, task will be marked as errored)
func (t *Task) Teardown(err error) {
	if err != nil {
		pkg.Logger.Printf("task %s: failed with error: %s", t.ID, err)
		t.Err = err
		t.Status = StatusError
	} else {
		t.Status = StatusReady
	}

	t.StopAt = time.Now().Unix()
}
