package fetcher

import (
	"os"
	"strings"
)

type file struct {
}

func NewFile() Fetcher {
	return file{}
}

func (file) Check(url string) bool {
	if !strings.HasPrefix(url, "file://") {
		return false
	}
	file := strings.TrimPrefix(url, "file://")
	return checkFileExist(file)
}

func (file) Fetch(url, _ string) (string, *SubInfo, error) {
	filePath := strings.TrimPrefix(url, "file://")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}
	return string(content), nil, nil
}

func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
