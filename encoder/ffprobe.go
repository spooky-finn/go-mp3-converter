package encoder

import (
	"os/exec"
	"strconv"
	"strings"

	"3205.team/go-mp3-converter/cfg"
)

// VideoDuration returns the duration of the video in seconds
func VideoDuration(path string) float64 {
	// Command to run ffprobe and get the duration of the video
	cmd := exec.Command(cfg.AppConfig.FfprobebinPath, "-i", path, "-show_entries", "format=duration", "-of", "csv=p=0")

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	r, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		panic(err)
	}
	return r
}
