package main

import (
	"is4/internal"
)

const (
	docIdsFilename    = "work/docks.gob"
	docsIndexFilename = "work/index.gob"
	trigramsFilename  = "work/trigrams.gob"
)

func getMainIndexer(args internal.Args) internal.Indexer[internal.Args] {
	return &internal.MainIndexer{
		Normalizer: internal.LowerCaseNormalizer{},
		Tokenizer:  internal.DefaultTokenizer{},
		ReaderProvider: &internal.ReaderProvider{
			ReadersMapper: &internal.ExtToReaderMapper{},
			DefaultReader: &internal.DefaultFileReader{},
			ExtProvider:   &internal.ExtByPathProvider{},
		},
		DocsProvider: &internal.DocsProvider{
			Provider: &internal.DocsByExtsProvider{
				Exts: args.Exts,
			},
		},
		DocsIndexer: &internal.DocsIndexer{
			IndexWriter: &internal.GobWriter[internal.DocsIndex]{
				FileName: docsIndexFilename,
			},
		},
		DocIdsWriter: &internal.GobWriter[internal.DocIds]{
			FileName: docIdsFilename,
		},
		TermIndexer: &TrigramIndexer{
			IndexWriter: &internal.GobWriter[TrigramIndex]{
				FileName: trigramsFilename,
			},
		},
	}
}

func getArgsSupplier() internal.ArgsSupplier {
	return internal.ArgsSupplier{
		Validator: internal.ArgsValidator{},
	}
}

func getMainSearcher() internal.Searcher[string, []string] {
	return internal.DocSearcher{
		IndexReader: &internal.GobReader[internal.DocsIndex]{
			FileName: docsIndexFilename,
		},
		DocsReader: &internal.GobReader[internal.DocIds]{
			FileName: docIdsFilename,
		},
		TermSearcher: TermSearcher{
			IndexReader: &internal.GobReader[TrigramIndex]{
				FileName: trigramsFilename,
			},
			PostFilter: &ByTemplateFiler{},
		},
	}

}
