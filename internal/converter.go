package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"

	"3205.team/go-mp3-converter/encoder"
	"3205.team/go-mp3-converter/helpers"
)

var ErrConverter = errors.New("converter error: ")

type ConvertEntity struct {
	filename      string
	DownloadURL   string
	outputFile    string
	inputFile     string
	inputTempFile string
	prog          *helpers.Progress
}

func NewConvertEntity(url string, prog *helpers.Progress) *ConvertEntity {
	filename := getFilename(url)
	return &ConvertEntity{
		filename:      filename,
		DownloadURL:   url,
		inputFile:     path.Join(TEMP_DIR, fmt.Sprintf("%s.mp4", filename)),
		inputTempFile: path.Join(TEMP_DIR, fmt.Sprintf("%s.mp4.temp", filename)),
		outputFile:    path.Join(TEMP_DIR, fmt.Sprintf("%s.mp3", filename)),
		prog:          prog,
	}
}

// Save file from url to temp dir and convert it to mp3
func (c *ConvertEntity) ConvertToMP3() (string, error) {
	err := SaveFileFromURL(c)
	if err != nil {
		helpers.Logger.Println(err)
		return "", errors.Join(ErrConverter, err)
	}

	ffmpeg := encoder.NewFfmpeg(c.inputFile, c.outputFile, c.prog)
	if err := ffmpeg.Run(); err != nil {
		helpers.Logger.Println(err)
		return "", errors.Join(ErrConverter, err)

	}

	if err := os.Remove(c.inputFile); err != nil {
		helpers.Logger.Println(err)
		return "", errors.Join(ErrConverter, err)
	}

	println("closing prog chan")
	close(c.prog.Ch)
	return c.outputFile, nil
}

func getFilename(url string) string {
	hashFunc := sha1.New()
	hashFunc.Write([]byte(url))
	return hex.EncodeToString(hashFunc.Sum(nil)[:12])
}
