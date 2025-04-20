package indexers

import (
	"is5/internal/core"
	"is5/internal/logger"
	"sort"
)

type BlockIndexer struct {
	MaxMapLength int
	DocSupplier  core.Supplier[core.Doc]
	TxtReader    core.DocReader
	Tokenizer    core.Tokenizer
	Normalizer   core.Normalizer
}

func (bi *BlockIndexer) Index() (*[]core.TermToDocIds, error) {
	logger.Log.Info("Indexing block")
	termsMap := make(map[string][]int)

	for {
		if len(termsMap) >= bi.MaxMapLength {
			break
		}
		doc, err := bi.DocSupplier.Supply()
		if err != nil {
			return nil, err
		}
		logger.Log.Debug("Indexing doc: ", doc.Path, " ", doc.Id)
		path := doc.Path
		id := doc.Id
		lines, err := bi.TxtReader.Read(path)
		if err != nil {
			return nil, err
		}
		for line := range lines {
			tokens := bi.Tokenizer.Tokenize(line)
			terms := bi.Normalizer.Normalize(tokens)
			for _, term := range terms {
				if termsMap[term] == nil {
					termsMap[term] = make([]int, 0)
				}
				if len(termsMap[term]) > 0 {
					if termsMap[term][len(termsMap[term])-1] < id {
						termsMap[term] = append(termsMap[term], id)
					}
				} else {
					termsMap[term] = append(termsMap[term], id)
				}

			}

		}
	}
	termsToDocIds := make([]core.TermToDocIds, len(termsMap))
	i := 0
	for term, docIds := range termsMap {
		termsToDocIds[i] = core.TermToDocIds{
			Term:   term,
			DocIds: docIds,
		}
		i++
	}
	sort.Slice(termsToDocIds, func(i, j int) bool {
		return termsToDocIds[i].Term < termsToDocIds[j].Term
	})
	return &termsToDocIds, nil
}
