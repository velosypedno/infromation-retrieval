package main

import (
	"encoding/gob"
	"is3/internal"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type TermPos struct {
	term string
	doc  int
}
type TermIndex struct {
	Term string
	Docs []int
}

type PositionalIndex struct {
	InvIndex []TermIndex
	Docs     map[int]string
}

func mapStep(pathIndexes map[string]int, termsPos *[]TermPos) error {
	for path, doc := range pathIndexes {
		var reader internal.FileReader
		if filepath.Ext(path) == ".fb2" {
			reader = internal.XMLFileReader{}
		} else {
			reader = internal.DefaultFileReader{}
		}
		lines, err := reader.Read(path)
		if err != nil {
			return err
		}
		tokenizer := internal.DefaultTokenizer{}
		normalizer := internal.LowerCaseNormalizer{}
		for line := range lines {
			tokens := tokenizer.Tokenize(line)
			terms := normalizer.Normalize(tokens)

			for i, term := range terms {
				if i == 0 {
					termPos := TermPos{term: term, doc: doc}
					*termsPos = append(*termsPos, termPos)
				} else {
					termPos := TermPos{term: term, doc: doc}
					*termsPos = append(*termsPos, termPos)
					termPos = TermPos{term: terms[i-1] + " " + term, doc: doc}
					*termsPos = append(*termsPos, termPos)
				}

			}
		}
	}
	return nil
}

func reduceStep(termsPos []TermPos, invIndex *[]TermIndex) error {
	invIndexMap := make(map[string]map[int]struct{})
	for _, termPos := range termsPos {
		if _, ok := invIndexMap[termPos.term]; !ok {
			invIndexMap[termPos.term] = make(map[int]struct{}, 0)
		}
		invIndexMap[termPos.term][termPos.doc] = struct{}{}
	}

	for term, docsMap := range invIndexMap {
		docs := make([]int, 0)
		for doc := range docsMap {
			docs = append(docs, doc)
		}
		termIndex := TermIndex{Term: term, Docs: docs}
		sort.Slice(termIndex.Docs, func(i, j int) bool {
			return termIndex.Docs[i] < termIndex.Docs[j]
		})
		*invIndex = append(*invIndex, termIndex)
	}
	sort.Slice(*invIndex, func(i, j int) bool {
		return (*invIndex)[i].Term < (*invIndex)[j].Term
	})
	return nil
}

func indexFiles(args internal.Args) {
	files, err := internal.GetDocsFromDirs(args.Dirs, args.Exts)
	if err != nil {
		log.Fatal(err)
	}

	pathToIndex := make(map[string]int)
	indexToPath := make(map[int]string)
	for index, path := range files {
		pathToIndex[path] = index
		indexToPath[index] = path
	}

	termsPos := make([]TermPos, 0)
	err = mapStep(pathToIndex, &termsPos)
	if err != nil {
		log.Fatal(err)
	}

	invIndex := []TermIndex{}
	err = reduceStep(termsPos, &invIndex)
	if err != nil {
		log.Fatal(err)
	}
	posIndex := PositionalIndex{InvIndex: invIndex, Docs: indexToPath}
	saveIndex(&posIndex)
}

func saveIndex(posIndex *PositionalIndex) {
	f, err := os.Create("biwordindex.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(posIndex)
	if err != nil {
		log.Fatal(err)
	}
}

func loadIndex() PositionalIndex {
	f, err := os.Open("biwordindex.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var posIndex PositionalIndex
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&posIndex)
	if err != nil {
		log.Fatal(err)
	}
	return posIndex
}

func binarySearch(term string, posIndex *PositionalIndex) int {
	start := 0
	end := len(posIndex.InvIndex) - 1
	for start <= end {
		mid := (start + end) / 2
		if posIndex.InvIndex[mid].Term < term {
			start = mid + 1
		} else if posIndex.InvIndex[mid].Term > term {
			end = mid - 1
		} else {
			return mid
		}
	}
	return -1
}
