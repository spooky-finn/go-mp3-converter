package internal

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"3205.team/go-mp3-converter/helpers"
)

const (
	TEMP_DIR      = "./temp"
	timerInterval = 1 * time.Second
)

var (
	ErrNoContentDisposition = errors.New("no content-disposition header")
	ErrPullFile             = errors.New("can't get file")
	ErrSaveToStorage        = errors.New("can't save file to storage")
	ErrNetDownloadFile      = errors.New("network failure while downloading file")
)

// ProgResponse counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type ProgResponse struct {
	downloaded float64
	total      float64
	progress   *helpers.Progress
	ticker     *time.Ticker
}

func (d *ProgResponse) Write(p []byte) (int, error) {
	n := len(p)
	d.downloaded += float64(n)
	select {
	case <-d.ticker.C:
		// send progress to channel every second\
		d.progress.Send(helpers.DownloadStage, int(d.downloaded/d.total*100))
	default:
	}
	return n, nil
}

func SaveFileFromURL(entity *ConvertEntity) error {
	out, err := os.Create(entity.inputTempFile)
	if err != nil {
		helpers.Logger.Println(err)
		return err
	}

	resp, err := http.Get(entity.DownloadURL)
	if err != nil {
		helpers.Logger.Println(ErrPullFile)
		return ErrPullFile
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		helpers.Logger.Println(ErrNetDownloadFile)
		return ErrNetDownloadFile
	}

	progResponse := &ProgResponse{
		total:    float64(resp.ContentLength),
		progress: entity.prog,
		ticker:   time.NewTicker(timerInterval),
	}
	defer progResponse.ticker.Stop()
	teeReader := io.TeeReader(resp.Body, progResponse)
	if _, err := io.Copy(out, teeReader); err != nil {
		helpers.Logger.Println(ErrSaveToStorage)
		return ErrSaveToStorage
	}

	if err := out.Close(); err != nil {
		return err
	}

	if err = os.Rename(entity.inputTempFile, entity.inputFile); err != nil {
		return err
	}

	return nil
}