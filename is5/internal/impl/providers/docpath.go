package providers

import (
	"is5/internal/logger"
	"os"
	"path/filepath"
)

type DocPathsByDirProvider struct{}

func (s *DocPathsByDirProvider) Provide(dir string) ([]string, error) {
	logger.Log.Info("Provide doc paths from dir ", dir)
	docPaths := make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			logger.Log.Debug("Found dir ", filepath.Join(dir, entry.Name()))
			paths, err := s.Provide(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			docPaths = append(docPaths, paths...)
		} else {
			logger.Log.Debug("Found file ", filepath.Join(dir, entry.Name()))
			docPaths = append(docPaths, filepath.Join(dir, entry.Name()))
		}
	}
	logger.Log.Info("No more docs in dir ", dir)
	logger.Log.Info("Found ", len(docPaths), " docs")
	return docPaths, nil
}
