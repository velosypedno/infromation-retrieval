package main

import (
	"log"
	"path/filepath"
)

type TermPosition struct {
	term      string
	fileIndex int
}

func getTerms(pathIndexes map[string]int, termsSet map[TermPosition]struct{}) error {
	for path := range pathIndexes {
		var reader FileReader
		if filepath.Ext(path) == ".fb2" {
			reader = XMLFileReader{}
		} else {
			reader = DefaultFileReader{}
		}
		lines, err := reader.Read(path)
		if err != nil {
			return err
		}
		tokenizer := DefaultTokenizer{}
		normalizer := LowerCaseNormalizer{}
		for line := range lines {
			tokens := tokenizer.Tokenize(line)
			terms := normalizer.Normalize(tokens)
			for _, term := range terms {
				termPosition := TermPosition{term: term, fileIndex: pathIndexes[path]}
				termsSet[termPosition] = struct{}{}
			}
		}
	}
	return nil
}

func indexFiles(args Args) {
	files, err := getFilesFromDirectories(args.dirs, args.exts)
	if err != nil {
		log.Fatal(err)
	}

	pathToIndex := make(map[string]int)
	indexToPath := make(map[int]string)
	for index, path := range files {
		pathToIndex[path] = index
		indexToPath[index] = path
	}

	termsSet := make(map[TermPosition]struct{})
	err = getTerms(pathToIndex, termsSet)
	if err != nil {
		log.Fatal(err)
	}

	err = buildMatrix(termsSet, indexToPath)
	if err != nil {
		log.Fatal(err)
	}

	err = buildIndex(termsSet, indexToPath)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := parseArgs()
	err := validateArgs(args)
	if err != nil {
		log.Fatal(err)
	}
	indexFiles(args)
	searchByMatrix(args.query)
	searchByIndex(args.query)

}
