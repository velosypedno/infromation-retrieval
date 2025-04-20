package suppliers

import (
	"is5/internal/core"
	"is5/internal/logger"
)

type DocPathsSupplier struct {
	Dirs                  []string
	DocPathsByDirProvider core.Provider[string, []string]
}

func (s *DocPathsSupplier) Supply() ([]string, error) {
	logger.Log.Info("Supplying doc paths")
	docPaths := make([]string, 0)
	for _, dir := range s.Dirs {
		paths, err := s.DocPathsByDirProvider.Provide(dir)
		if err != nil {
			return nil, err
		}
		docPaths = append(docPaths, paths...)
	}
	logger.Log.Info("Found ", len(docPaths), " docs")
	return docPaths, nil

}
