package entity

import "3205.team/go-mp3-converter/pkg"

type Status string

const (
	StatusNew      Status = "new"
	StatusReady    Status = "ready"
	StatusError    Status = "error"
	StatusProgress Status = "in_progress"
)

// a main entity that represents a task for convertation to mp3
type Task struct {
	ID string `json:"filename"`
	// a url to the video (e.g. to the youtube video)
	OriginalURL string `json:"originalUrl"`
	// a url issued by Link microservice to download the video
	DownloadURL string        `json:"downloadUrl"`
	Status      Status        `json:"status,omitempty"`
	Progress    *pkg.Progress `json:"-"`
	Duration    float64       `json:"duration"`
	Thumb       string        `json:"thumb"`

	PushedAt int64 `json:"pushedAt"`
	StartAt  int64 `json:"startAt"`
	StopAt   int64 `json:"stopAt"`
}

type NewTaskParams struct {
	OriginalURL string
	DownloadURL string
}

func NewTask(params NewTaskParams) *Task {
	// if params.OriginalURL == "" || params.DownloadURL == "" {
	// 	panic("empty params")
	// }

	ID := pkg.Hash(params.OriginalURL)
	return &Task{
		ID:          ID,
		Progress:    pkg.NewProg(),
		OriginalURL: params.OriginalURL,
		Status:      StatusNew,
		DownloadURL: params.DownloadURL,
	}
}
