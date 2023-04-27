package encoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"3205.team/go-mp3-converter/helpers"
)

const FFMPEG_BIN = "./bin/ffmpeg"

var Errffmpeg = errors.New("can't convert file to mp3")

type ffmpeg struct {
	inputFile  string
	outputFile string
	videoDur   float64
	prog       *helpers.Progress
}

func NewFfmpeg(inputFile, outputFile string, prog *helpers.Progress) *ffmpeg {
	return &ffmpeg{
		inputFile:  inputFile,
		outputFile: outputFile,
		videoDur:   VideoDuration(inputFile),
		prog:       prog,
	}
}

func (f *ffmpeg) Run() error {
	args := []string{"-i",
		f.inputFile,
		"-y", "-map", "0:a:0", "-acodec", "libmp3lame", "-q:a", "1", "-f", "mp3", "-progress", "pipe:1",
		f.outputFile}

	cmd := exec.Command(FFMPEG_BIN, args...)
	ffmpegStdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("ffmpeg err: %v", err)
		return Errffmpeg
	}

	f.updateProgress(ffmpegStdout)
	cmd.Wait()
	return nil
}

func (f *ffmpeg) updateProgress(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// out_time_us is a already processed time of mp3 in microseconds
		if strings.Contains(line, "out_time_us=") {
			t, err := strconv.ParseFloat(strings.Split(line, "=")[1], 64)
			if err != nil {
				panic(err)
			}
			f.prog.Send(helpers.ConvertStage, int((t/1000000/f.videoDur)*100))
		}
	}
}
