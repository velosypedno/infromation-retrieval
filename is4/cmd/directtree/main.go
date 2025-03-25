package main

import (
	"is4/internal"
	"log"
	"time"
)

func main() {
	argsSupplier := getArgsSupplier()
	args, err := argsSupplier.Supply()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(args)
	switch args.Mode {
	case internal.Index:
		index(args)
	case internal.Search:
		search(args)
	case internal.IndexSearch:
		index(args)
		log.Println("---------------------------------------")
		search(args)
	}

}

func index(args internal.Args) {
	start := time.Now()

	log.Println("Indexing...")
	indexer := getMainIndexer(args)
	err := indexer.Index(args)
	if err != nil {
		log.Fatal(err)
	}

	duration := time.Since(start)
	log.Println("Indexing finished in", duration)
}

func search(args internal.Args) {
	start := time.Now()
	log.Println("Searching...")
	searcher := getMainSearcher()
	terms, err := searcher.Search(args.Query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Results:")
	for i, term := range terms {
		log.Println(i, term)
	}

	duration := time.Since(start)
	log.Println("Searching finished in", duration)
}
