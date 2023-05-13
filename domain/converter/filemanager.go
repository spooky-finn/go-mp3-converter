package converter

import (
	"fmt"
	"io"
	"os"
	"path"

	"3205.team/go-mp3-converter/pkg"
)

type FilePaths struct {
	// an mp4 file that will be converted to mp3
	Original string
	// a filepath have that name while it is downloading
	OriginalTemp string
	// a result mp3 file
	Output string
}

type FileManager struct {
	FilePaths
}

func NewFileManager(tempDir string, filename string) *FileManager {
	return &FileManager{
		FilePaths: FilePaths{
			Original:     path.Join(tempDir, fmt.Sprintf("%s.mp4", filename)),
			OriginalTemp: path.Join(tempDir, fmt.Sprintf("%s.mp4.temp", filename)),
			Output:       path.Join(tempDir, fmt.Sprintf("%s.mp3", filename)),
		},
	}
}

// create a file with a name OriginalTemp and copy copy content from io.Reader
func (fm *FileManager) CopyToTempFile(r io.Reader) error {
	tempFile, err := os.Create(fm.OriginalTemp)
	if err != nil {
		pkg.Logger.Fatalln(err)
		return err
	}

	if _, err := io.Copy(tempFile, r); err != nil {
		pkg.Logger.Fatalln(err)
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	fm.renameTempFileToOriginal()
	return nil
}

func (fm *FileManager) renameTempFileToOriginal() error {
	err := os.Rename(fm.OriginalTemp, fm.Original)
	if err != nil {
		pkg.Logger.Fatalln(err)
		return err
	}
	return nil
}

func (fm *FileManager) RemoveOriginalFile() error {
	err := os.Remove(fm.Original)
	if err != nil {
		pkg.Logger.Fatalln(err)
		return err
	}

	return nil
}

func (fm *FileManager) RemoveOutputFile() error {
	err := os.Remove(fm.Output)
	if err != nil {
		pkg.Logger.Fatalln(err)
		return err
	}

	return nil
}
