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
	pos  int
}
type TermIndex struct {
	Term string
	Docs []struct {
		Doc  int
		Poss []int
	}
}

type PositionalIndex struct {
	InvIndex []TermIndex
	Docs     map[int]string
}

func mapStep(pathIndexes map[string]int, termsPos *[]TermPos) error {
	var pos int
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
		pos = 0
		tokenizer := internal.DefaultTokenizer{}
		normalizer := internal.LowerCaseNormalizer{}
		for line := range lines {
			tokens := tokenizer.Tokenize(line)
			terms := normalizer.Normalize(tokens)
			for _, term := range terms {
				termPos := TermPos{term: term, doc: doc, pos: pos}
				pos++
				*termsPos = append(*termsPos, termPos)
			}
		}
	}
	return nil
}

func reduceStep(termsPos []TermPos, invIndex *[]TermIndex) error {
	invIndexMap := make(map[string]map[int][]int)
	for _, termPos := range termsPos {
		if _, ok := invIndexMap[termPos.term]; !ok {
			invIndexMap[termPos.term] = make(map[int][]int)
		}
		invIndexMap[termPos.term][termPos.doc] = append(invIndexMap[termPos.term][termPos.doc], termPos.pos)
	}

	for term, docs := range invIndexMap {
		termIndex := TermIndex{Term: term}
		for doc, positions := range docs {
			termIndex.Docs = append(termIndex.Docs, struct {
				Doc  int
				Poss []int
			}{Doc: doc, Poss: positions})
		}
		sort.Slice(termIndex.Docs, func(i, j int) bool {
			return termIndex.Docs[i].Doc < termIndex.Docs[j].Doc
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
	f, err := os.Create("posindex.gob")
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
	f, err := os.Open("posindex.gob")
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
