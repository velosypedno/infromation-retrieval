package usecases

import (
	"is5/internal/core"
	"is5/internal/logger"
	"sync"
)

type SPIMI struct {
	IndexersAmount     int
	DocsTxtFilePath    string
	BlockIndexer       core.Indexer
	IndexWriter        core.Writer[*[]core.TermToDocIds]
	DocProducer        core.Producer
	IndexPairsProducer core.Producer
	MergeWorker        core.Worker
}

func (s *SPIMI) BuildIndex() error {
	logger.Log.Info("Producing doc paths to chan")
	err := s.DocProducer.Produce()
	if err != nil {
		return err
	}
	var indexersWG sync.WaitGroup
	for i := 0; i < s.IndexersAmount; i++ {
		indexersWG.Add(1)
		logger.Log.Info("Starting indexer with id ", i)
		go func() {
			id := i
			for {
				termsToDocIds, err := s.BlockIndexer.Index()
				if err != nil {
					indexersWG.Done()
					break
				}
				err = s.IndexWriter.Write(termsToDocIds)
				if err != nil {
					logger.Log.Error(err)
				}
			}
			logger.Log.Info("Indexer with id ", id, " finished")
		}()
	}
	logger.Log.Info("Producing index pairs to chan")
	err = s.IndexPairsProducer.Produce()
	if err != nil {
		return err
	}

	err = s.MergeWorker.Work()
	if err != nil {
		return err
	}

	indexersWG.Wait()
	return nil
}
