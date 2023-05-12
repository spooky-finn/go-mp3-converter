package converter

import (
	"errors"
	"time"

	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain/encoder"
	"3205.team/go-mp3-converter/entity"
	"3205.team/go-mp3-converter/pkg"
)

var ErrConverter = errors.New("converter error: ")

// type ConvertEntity struct {
// 	Filename      string `json:"filename"`
// 	DownloadURL   string
// 	Status        Status  `json:"status"`
// 	StartAt       int64   `json:"startAt"`
// 	StopAt        int64   `json:"stopAt"`
// 	Duration      float64 `json:"duration"`
// 	Thumb         string  `json:"thumb"`
// 	outputFile    string
// 	inputFile     string
// 	inputTempFile string
// 	prog          *helpers.Progress
// }

// func NewConvertEntity(url string, prog *helpers.Progress) *ConvertEntity {
// 	filename := helpers.Hash(url)
// 	dir := cfg.AppConfig.TempDir
// 	return &ConvertEntity{
// 		Filename:      filename,
// 		DownloadURL:   url,
// 		Status:        StatusNew,
// 		inputFile:     path.Join(dir, fmt.Sprintf("%s.mp4", filename)),
// 		inputTempFile: path.Join(dir, fmt.Sprintf("%s.mp4.temp", filename)),
// 		outputFile:    path.Join(dir, fmt.Sprintf("%s.mp3", filename)),
// 		prog:          prog,
// 	}
// }

type MP3ConverterUseCase struct {
}

// Save file from url to temp dir and convert it to mp3
func Run(t *entity.Task) error {
	t.Status = entity.StatusProgress
	t.StartAt = time.Now().Unix()

	tempDir := cfg.AppConfig.TempDir
	fm := NewFileManager(tempDir, t.ID)

	err := SaveFileFromURL(t, fm)
	if err != nil {
		pkg.Logger.Println(err)
		return errors.Join(ErrConverter, err)
	}
	t.Duration = encoder.GetVideoDuration(fm.Original)

	ffmpeg := &encoder.Ffmpeg{
		InputFile:  fm.Original,
		OutputFile: fm.Output,
		VideoDur:   t.Duration,
		Prog:       t.Progress,
	}
	if err := ffmpeg.Run(); err != nil {
		pkg.Logger.Println(err)
		return errors.Join(ErrConverter, err)
	}

	fm.RemoveOriginalFile()

	t.StopAt = time.Now().Unix()
	t.Status = entity.StatusError
	close(t.Progress.Ch)
	return nil
}
