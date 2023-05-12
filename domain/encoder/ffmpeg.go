package encoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	cfg "3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/pkg"
)

var Errffmpeg = errors.New("can't convert file to mp3")

type Ffmpeg struct {
	InputFile  string
	OutputFile string
	VideoDur   float64
	Prog       *pkg.Progress
}

func (f *Ffmpeg) Run() error {
	args := []string{"-i",
		f.InputFile,
		"-y", "-map", "0:a:0", "-acodec", "libmp3lame", "-q:a", "1", "-f", "mp3", "-progress", "pipe:1",
		f.OutputFile}

	cmd := exec.Command(cfg.AppConfig.FfmpegbinPath, args...)
	ffmpegStdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("ffmpeg err: %v", err)
		return Errffmpeg
	}

	f.observeProgress(ffmpegStdout)
	cmd.Wait()
	return nil
}

func (f *Ffmpeg) observeProgress(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// out_time_us is a already processed time of mp3 in microseconds
		if strings.Contains(line, "out_time_us=") {
			t, err := strconv.ParseFloat(strings.Split(line, "=")[1], 64)
			if err != nil {
				panic(err)
			}
			f.Prog.Send(pkg.ConvertStage, int((t/1000000/f.VideoDur)*100))
		}
	}
}
