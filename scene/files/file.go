package files

import (
	"path/filepath"
)

func GetFilePath(fileName string) string {
	absPath, _ := filepath.Abs("./")
	joinedPath := filepath.Join(absPath, "files")

	return filepath.Join(joinedPath, filepath.Clean(fileName))
}
