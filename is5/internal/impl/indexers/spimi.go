package indexers

import (
	"is5/internal/core"
	"is5/internal/logger"
	"log"
	"sync"
)

type SPIMI struct {
	IndexersAmount     int
	DocsTxtFilePath    string
	DocsTxtFileCreator core.Creator[string]
	DocPathsSupplier   core.Supplier[[]string]
	DocsWriter         core.Writer[[]core.Doc]
	IndexWriter        core.Writer[*[]core.TermToDocIds]
	BlockIndexer       core.Indexer
	DocProducer        core.Producer
	IndexPairsProducer core.Producer
	MergeWorker        core.Worker
}

func (s *SPIMI) BuildIndex() error {
	logger.Log.Info("SPIMI algorithm is started")
	logger.Log.Info("Creating file for docs to index mapping")
	err := s.DocsTxtFileCreator.Create(s.DocsTxtFilePath)
	if err != nil {
		return err
	}

	docs := make([]core.Doc, 0)
	docPaths, err := s.DocPathsSupplier.Supply()
	if err != nil {
		return err
	}

	logger.Log.Info("Give id to each doc")
	for i, path := range docPaths {
		docs = append(docs, core.Doc{
			Path: path,
			Id:   i,
		})
	}
	logger.Log.Info("Max doc index: ", len(docs)-1)
	err = s.DocsWriter.Write(docs)
	if err != nil {
		return err
	}

	logger.Log.Info("Producing doc paths to chan")
	err = s.DocProducer.Produce()
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
					log.Println(err)
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
