package tempcleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"3205.team/go-mp3-converter/pkg"
)

type TempCleaner struct {
	tempFolder    string
	maxAge        time.Duration
	checkInterval time.Duration
}

var (
	instance *TempCleaner
)

func New(tempFolder string, maxAge time.Duration, checkInterval time.Duration) *TempCleaner {
	if instance != nil {
		panic("Another instance of TempCleaner already exists")
	}

	tc := &TempCleaner{
		maxAge:        maxAge,
		tempFolder:    tempFolder,
		checkInterval: checkInterval,
	}
	instance = tc
	tc.init()
	return tc
}

func (tc *TempCleaner) init() {
	go func() {
		for {
			tc.cleanTempFolder()
			time.Sleep(tc.checkInterval)
		}
	}()
}

func (tc *TempCleaner) cleanTempFolder() {
	err := filepath.Walk(tc.tempFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && tc.isFileOutdated(info.ModTime()) {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove file: %s", path)
			}

			return fmt.Errorf("removed file: %s", path)
		}

		return nil
	})

	if err != nil {
		pkg.Logger.Printf("Error while cleaning temp folder: %v\n", err)
	}
}

func (tc *TempCleaner) isFileOutdated(modTime time.Time) bool {
	return time.Since(modTime) > tc.maxAge
}
