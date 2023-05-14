package mp3converter

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain/encoder"
	"3205.team/go-mp3-converter/entity"
)

var ErrConverter = errors.New("converter error: ")

type TaskCacher interface {
	GetTask(ID string) *entity.Task
	SetTask(t *entity.Task)
}

type ConverterParams struct {
	DownloadURL string
	OriginalURL string
}

type MP3ConverterUseCase struct {
	cache TaskCacher
}

func New(cache TaskCacher) *MP3ConverterUseCase {
	return &MP3ConverterUseCase{
		cache: cache,
	}
}

// Starts convertation, but returns pointer of the task almost immediately
func (uc *MP3ConverterUseCase) StartConvertation(p ConverterParams) *entity.Task {
	t := entity.NewTask(entity.NewTaskParams{
		OriginalURL: p.OriginalURL,
		DownloadURL: p.DownloadURL,
	})
	t.Status = entity.StatusProgress
	t.StartAt = time.Now().Unix()

	cachedTask := uc.cache.GetTask(t.ID)
	if cachedTask != nil {
		return cachedTask
	}

	go uc.convert(t)
	return t
}

func (uc *MP3ConverterUseCase) convert(t *entity.Task) {
	tempDir := cfg.AppConfig.TempDir
	fm := NewFileManager(tempDir, t.ID)

	err := SaveFileFromURL(t, fm)
	if err != nil {
		t.Teardown(errors.Join(ErrConverter, err))
		return
	}
	t.Duration = encoder.GetVideoDuration(fm.Original)

	ffmpeg := &encoder.Ffmpeg{
		InputFile:  fm.Original,
		OutputFile: fm.Output,
		VideoDur:   t.Duration,
		Prog:       t.Progress,
	}
	if err := ffmpeg.Run(); err != nil {
		t.Teardown(errors.Join(ErrConverter, err))
		return
	}

	uc.cache.SetTask(t)
	fm.RemoveOriginalFile()
	t.Teardown(nil)
}
