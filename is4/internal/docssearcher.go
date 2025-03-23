package internal

import (
	"strings"
)

type DocSearcher struct {
	IndexReader  Reader[DocsIndex]
	DocsReader   Reader[DocIds]
	TermSearcher Searcher[string, []string]
}

func (d DocSearcher) Search(term string) ([]string, error) {
	terms, err := d.TermSearcher.Search(term)
	if err != nil {
		return nil, err
	}
	index, err := d.IndexReader.Read()
	if err != nil {
		return nil, err
	}
	docIds, err := d.DocsReader.Read()
	if err != nil {
		return nil, err
	}
	termDocsMap := make(map[string]string)

	for _, term := range terms {
		pos := binarySearch(term, index)
		if pos == -1 {
			continue
		}
		docs := []string{}
		for _, doc := range (*index)[pos].Docs {
			docs = append(docs, (*docIds)[doc])
		}
		termDocsMap[term] = strings.Join(docs, " ")

	}
	var result []string
	for term, docs := range termDocsMap {
		result = append(result, "Term: "+term+" Docs: ["+docs+"]")
	}
	return result, nil
}

func binarySearch(term string, index *DocsIndex) int {
	start := 0
	end := len(*index) - 1
	for start <= end {
		mid := (start + end) / 2
		if (*index)[mid].Term < term {
			start = mid + 1
		} else if (*index)[mid].Term > term {
			end = mid - 1
		} else {
			return mid
		}
	}
	return -1
}
