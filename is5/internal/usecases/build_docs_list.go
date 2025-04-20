package usecases

import (
	"is5/internal/core"
	"is5/internal/logger"
)

type DocsListBuilder struct {
	DocsTxtFilePath    string
	DocPathsSupplier   core.Supplier[[]string]
	DocsTxtFileCreator core.Creator[string]
	DocsWriter         core.Writer[[]core.Doc]
}

func (b *DocsListBuilder) BuildDocsList() error {
	logger.Log.Info("Creating file for docs to index mapping")
	err := b.DocsTxtFileCreator.Create(b.DocsTxtFilePath)
	if err != nil {
		return err
	}

	docs := make([]core.Doc, 0)
	docPaths, err := b.DocPathsSupplier.Supply()
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
	err = b.DocsWriter.Write(docs)
	if err != nil {
		return err
	}
	return nil
}
