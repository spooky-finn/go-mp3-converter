package mp3converter

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain"
	"3205.team/go-mp3-converter/domain/encoder"
	"3205.team/go-mp3-converter/entity"
)

var ErrConverter = errors.New("converter error: ")

type MP3Converter struct {
}

func New() *MP3Converter {
	return &MP3Converter{}
}

// Starts convertation, but returns pointer of the task almost immediately
func (uc *MP3Converter) StartConvertation(p entity.NewTaskParams, prog domain.Progress) *entity.Task {
	task := entity.NewTask(p)
	task.Status = entity.StatusInProgress
	task.StartAt = time.Now().Unix()

	go uc.convert(task, prog)
	return task
}

func (uc *MP3Converter) convert(t *entity.Task, prog domain.Progress) {
	tempDir := cfg.AppConfig.TempDir
	fm := NewFileManager(tempDir, t.ID)

	err := SaveFileFromURL(t.DownloadURL, fm, prog)
	if err != nil {
		t.Teardown(errors.Join(ErrConverter, err))
		return
	}
	t.Duration = encoder.GetVideoDuration(fm.Original)

	ffmpeg := &encoder.Ffmpeg{
		InputFile:  fm.Original,
		OutputFile: fm.Output,
		VideoDur:   t.Duration,
		Prog:       prog,
	}
	if err := ffmpeg.Run(); err != nil {
		t.Teardown(errors.Join(ErrConverter, err))
		return
	}

	t.Teardown(nil)

	fm.RemoveOriginalFile()
}
