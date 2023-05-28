package test

import (
	"os"
	"path"
	"time"

	"3205.team/go-mp3-converter/cfg"
)

type TempDirWatchr struct {
	Files []os.DirEntry
}

func NewTempDirWatcher() *TempDirWatchr {
	w := &TempDirWatchr{}
	go w.watch()
	return w
}

func (t *TempDirWatchr) IsEmpty() bool {
	return len(t.Files) == 0
}

func (t *TempDirWatchr) IsNotEmpty() bool {
	return len(t.Files) != 0
}

func (t *TempDirWatchr) HaveOnlyMp3files() bool {
	for _, each := range t.Files {
		if each.Name()[len(each.Name())-4:] != ".mp3" {
			return false
		}
	}
	return true
}

func (t *TempDirWatchr) IsLenOf(n int) bool {
	return len(t.Files) == n
}

func (t *TempDirWatchr) getFiles() (files []os.DirEntry) {
	tempDir := path.Join("../", cfg.AppConfig.TempDir)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		panic(err)
	}
	return files
}

func (t *TempDirWatchr) watch() {
	for {
		time.Sleep(1 * time.Second)
		t.Files = t.getFiles()
	}
}
