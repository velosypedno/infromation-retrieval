package workers

import (
	"is5/internal/core"
	"is5/internal/logger"
)

type MergeWorker struct {
	IndexFilesPairs            chan core.Pair[string]
	MergedIndexFilesToBeMerged chan string
	Merger                     core.Merger
}

func (m MergeWorker) Work() error {
	for indexFilesPair := range m.IndexFilesPairs {
		logger.Log.Info("Merging index files: ", indexFilesPair.First, " and ", indexFilesPair.Second)
		merged, err := m.Merger.Merge(indexFilesPair.First, indexFilesPair.Second)
		if err != nil {
			return err
		}
		logger.Log.Info("Merged index file: ", merged)
		logger.Log.Debug("Writing merged index file to channel: ", merged, ", channel: ", m.MergedIndexFilesToBeMerged)
		m.MergedIndexFilesToBeMerged <- merged
	}
	close(m.MergedIndexFilesToBeMerged)
	return nil
}
