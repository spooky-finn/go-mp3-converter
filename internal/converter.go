package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"3205.team/go-mp3-converter/encoder"
	"3205.team/go-mp3-converter/helpers"
)

type ConvertEntity struct {
	filename      string
	DownloadURL   string
	outputFile    string
	inputFile     string
	inputTempFile string
	prog          *helpers.Progress
}

func NewConvertEntity(url string, prog *helpers.Progress) *ConvertEntity {
	filename := fmt.Sprintf("%d", time.Now().UnixNano())
	return &ConvertEntity{
		filename:      filename,
		DownloadURL:   url,
		inputFile:     path.Join(TEMP_DIR, fmt.Sprintf("%s.mp4", filename)),
		inputTempFile: path.Join(TEMP_DIR, fmt.Sprintf("%s.mp4.temp", filename)),
		outputFile:    path.Join(TEMP_DIR, fmt.Sprintf("%s.mp3", filename)),
		prog:          prog,
	}
}

func (c *ConvertEntity) ConvertToMP3() (string, error) {

	err := DownloadFile(c)
	if err != nil {
		logger.Println(err)
		return "", err
	}

	ffmpeg := encoder.NewFfmpeg(c.inputFile, c.outputFile, c.prog)
	if err := ffmpeg.Run(); err != nil {
		logger.Println(err)
		return "", err

	}

	if err := os.Remove(c.inputFile); err != nil {
		logger.Println(err)
		return "", err
	}

	println("closing prog chan")
	close(c.prog.Ch)
	return c.outputFile, nil
}
