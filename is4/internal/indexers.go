package internal

import (
	"sort"
)

type DocsIndexer struct {
	IndexWriter Writer[DocsIndex]
}

func (d *DocsIndexer) Index(terms []TermDoc) error {
	docsIndexMap := make(map[string]map[int]struct{})

	for _, termDoc := range terms {
		if _, ok := docsIndexMap[termDoc.Term]; !ok {
			docsIndexMap[termDoc.Term] = make(map[int]struct{})
		}
		docsIndexMap[termDoc.Term][termDoc.Doc] = struct{}{}
	}
	var docsIndex DocsIndex
	for term, docsMap := range docsIndexMap {
		docs := make([]int, 0, len(docsMap))
		for doc := range docsMap {
			docs = append(docs, doc)
		}
		docsIndex = append(docsIndex, struct {
			Term string
			Docs []int
		}{Term: term, Docs: docs})
	}

	for _, termDocs := range docsIndex {
		sort.Slice(termDocs.Docs, func(i, j int) bool {
			return termDocs.Docs[i] < termDocs.Docs[j]
		})
	}
	sort.Slice(docsIndex, func(i, j int) bool {
		return docsIndex[i].Term < docsIndex[j].Term
	})
	err := d.IndexWriter.Write(&docsIndex)
	if err != nil {
		return err
	}
	return nil
}

type MainIndexer struct {
	Normalizer     Normalizer
	Tokenizer      Tokenizer
	DocsProvider   Provider[[]string, []string]
	ReaderProvider Provider[string, FileReader]
	DocsIndexer    Indexer[[]TermDoc]
	TermIndexer    Indexer[[]TermDoc]
	DocIdsWriter   Writer[DocIds]
}

func (di *MainIndexer) Index(args Args) error {
	collection := []TermDoc{}
	docs, err := di.DocsProvider.Provide(args.Dirs)
	if err != nil {
		return err
	}

	docToIndex := make(map[string]int)
	indexToDoc := make(DocIds)
	for i, doc := range docs {
		docToIndex[doc] = i
		indexToDoc[i] = doc
	}

	for _, doc := range docs {
		reader, err := di.ReaderProvider.Provide(doc)
		if err != nil {
			return err
		}
		stream, err := reader.Read(doc)
		if err != nil {
			return err
		}
		for line := range stream {
			tokens := di.Tokenizer.Tokenize(line)
			terms := di.Normalizer.Normalize(tokens)
			for _, term := range terms {
				collection = append(collection, TermDoc{
					Term: term,
					Doc:  docToIndex[doc],
				})
			}
		}
	}
	err = di.DocsIndexer.Index(collection)
	if err != nil {
		return err
	}
	err = di.TermIndexer.Index(collection)
	if err != nil {
		return err
	}
	err = di.DocIdsWriter.Write(&indexToDoc)
	if err != nil {
		return err
	}
	return nil
}
