package util

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type FileUtil struct {
}

func NewFileUtil() *FileUtil {
	return &FileUtil{}
}
func (util *FileUtil) RemoveAllFilesInFolder(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		if file.IsDir() {
			err := util.RemoveAllFilesInFolder(filePath)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (util *FileUtil) Download(url string, path string) (err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	return err
}
