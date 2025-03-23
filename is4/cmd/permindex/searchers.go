package main

import (
	"is4/internal"
	"strings"
)

type TermSearcher struct {
	IndexReader    internal.Reader[PermIndex]
	PrefixProvider internal.Provider[string, string]
}

func (t *TermSearcher) Search(template string) ([]string, error) {
	index, err := t.IndexReader.Read()
	if err != nil {
		return nil, err
	}
	prefix, err := t.PrefixProvider.Provide(template)
	if err != nil {
		return nil, err
	}

	terms := []string{}
	pos := len(*index) - 1
	endPos := pos
	for prefix < (*index)[pos].Perm {
		endPos = pos
		pos = pos / 2
	}
	for strings.HasPrefix((*index)[pos].Perm, prefix) {
		pos = pos / 2
	}

	for pos < endPos {
		if strings.HasPrefix((*index)[pos].Perm, prefix) {
			terms = append(terms, (*index)[pos].Term)
		}
		pos++
	}
	return terms, nil
}
