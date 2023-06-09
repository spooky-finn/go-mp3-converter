package mp3converter

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"3205.team/go-mp3-converter/application/progress"
	"3205.team/go-mp3-converter/cfg"
	"3205.team/go-mp3-converter/domain"
)

// var (
// 	ErrNoContentDisposition = errors.New("no content-disposition header")
// 	ErrPullFile             = errors.New("can't get file")
// 	ErrSaveToStorage        = errors.New("can't save file to storage")
// 	ErrNetDownloadFile      = errors.New("response status is not ok")
// )

// DownloadProgWatcher counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type DownloadProgWatcher struct {
	downloaded float64
	total      float64
	progress   domain.Progress
	ticker     *time.Ticker
}

func NewDownloadProgWatcher(totalBytes int64, progress domain.Progress) *DownloadProgWatcher {
	return &DownloadProgWatcher{
		total:    float64(totalBytes),
		progress: progress,
		ticker:   time.NewTicker(cfg.AppConfig.ProgressInterval),
	}
}

func (d *DownloadProgWatcher) Write(p []byte) (int, error) {
	n := len(p)
	d.downloaded += float64(n)
	select {
	case <-d.ticker.C:
		// send progress to channel every second\
		d.progress.Send(progress.DownloadStage, int(d.downloaded/d.total*100))
	default:
	}
	return n, nil
}

func SaveFileFromURL(downloadURL string, fileManager *FileManager, prog domain.Progress) error {
	client := http.Client{
		Timeout: 15 * time.Minute,
	}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status is not ok: %d", resp.StatusCode)
	}

	progResponse := NewDownloadProgWatcher(resp.ContentLength, prog)
	defer progResponse.ticker.Stop()

	teeReader := io.TeeReader(resp.Body, progResponse)
	if err := fileManager.CopyToTempFile(teeReader); err != nil {
		return err
	}
	return nil
}
