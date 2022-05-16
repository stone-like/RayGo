package files

import (
	"path/filepath"
	"runtime"
)

func getFileDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := filepath.Join(filepath.Dir(b))

	return d
}

var fileDir = getFileDir()

func GetFilePath(fileName string) string {

	return filepath.Join(fileDir, filepath.Clean(fileName))
}
