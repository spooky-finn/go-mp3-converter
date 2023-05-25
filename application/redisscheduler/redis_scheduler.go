package redisscheduler

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/application/progress"
	"3205.team/go-mp3-converter/domain/mp3converter"
	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"
	"github.com/redis/go-redis/v9"
)

type OutgoingTaskMeta struct {
	Title    string `json:"title"`
	Source   string `json:"source"`
	Duration string `json:"duration"`
	Tags     string `json:"tags"`
}

type OutgoingTask struct {
	Progress    int              `json:"progress"`
	OriginalURL string           `json:"originalUrl"`
	TaskName    string           `json:"taskName"`
	Status      entity.Status    `json:"status"`
	Origin      string           `json:"origin"`
	PushedAt    int64            `json:"pushedAt"`
	Instance    string           `json:"instance"`
	Meta        OutgoingTaskMeta `json:"meta"`
	Thumb       string           `json:"thumb"`
	FileSize    int64            `json:"fileSize"`
	StartAt     int64            `json:"startAt"`
	Duration    int              `json:"duration"`
	Filename    string           `json:"filename"`
	StopAt      int64            `json:"stopAt"`
	Errormsg    string           `json:"errorMsg"`
}

type RedisScheduler struct {
	taskhandler *Handler
	pooler      *Pooler
}

func NewRedisScheduler(redisclient *redis.Client, mp3converter *mp3converter.MP3Converter) *RedisScheduler {
	rs := &RedisScheduler{
		taskhandler: NewHandler(mp3converter),
		pooler:      NewPooler(redisclient),
	}
	go rs.init()
	return rs
}

func (r *RedisScheduler) init() {
	for request := range r.pooler.Ch {
		go r.handleIncomingRequest(request)
	}
}

func (r *RedisScheduler) handleIncomingRequest(request *IncomingTask) {
	prog := progress.New()
	result := r.taskhandler.Handle(request, prog)

	if result.Err != nil {
		if errors.Is(result.Err, ErrQueueTimeout) {
			pkg.Logger.Println("queue timeout elapsed")
		} else {
			panic(result.Err)
		}
	}

	for prog_event := range prog.Ch {
		println("progress in redis scheduler")
		outgoingTask := convertToOutgoingTask(result, request, prog_event)

		if err := r.pooler.Push(outgoingTask, 1*time.Hour); err != nil {
			panic(err)
		}
		continue
	}

}

func convertToOutgoingTask(result *HadlerResult, request *IncomingTask, prog_event progress.ProgressEventPayload) *OutgoingTask {
	task := result.Task
	linkResponse := result.LinkResponse

	// error
	errorMsg := ""
	if task.Err != nil {
		errorMsg = task.Err.Error()
	}

	return &OutgoingTask{
		Progress:    prog_event.AggragetedProgress(),
		OriginalURL: request.OriginalURL,
		TaskName:    request.TaskName,
		Origin:      request.Origin,
		PushedAt:    request.PushedAt,
		Instance:    "instance",
		Meta: OutgoingTaskMeta{
			Title:  linkResponse.Meta.Title,
			Source: linkResponse.Meta.Source,
		},
		Thumb:    linkResponse.Thumb,
		Status:   task.Status,
		FileSize: task.FileSize,
		StartAt:  task.StartAt,
		StopAt:   task.StopAt,
		Duration: int(task.Duration),
		Filename: task.ID + ".mp3",
		Errormsg: errorMsg,
	}
}
