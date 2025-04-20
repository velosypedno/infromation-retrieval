package creators

import (
	"is5/internal/logger"
	"os"
	"path/filepath"
)

type FileByPathCreator struct{}

func (f FileByPathCreator) Create(path string) error {
	logger.Log.Debug("Creating file by path")
	dir := filepath.Dir(path)
	os.Mkdir(dir, os.ModePerm)
	logger.Log.Debug("Creating dir for file:", dir)
	file, err := os.Create(path)
	if err != nil {
		logger.Log.Error("Unexpected error creating file: ", err)
		return err
	}
	logger.Log.Debug("File - ", path, " was created")
	defer func() {
		file.Close()
		logger.Log.Debug("File - ", path, " was closed")
	}()
	return nil
}
