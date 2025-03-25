package main

import (
	"is4/internal"
)

type TreeIndex struct {
	Terms []string
}

type TreeIndexer struct {
	IndexWriter internal.Writer[TreeIndex]
}

func (ti *TreeIndexer) Index(terms []internal.TermDoc) error {
	treeIndex := TreeIndex{}
	for _, term := range terms {
		treeIndex.Terms = append(treeIndex.Terms, term.Term)
	}
	return ti.IndexWriter.Write(&treeIndex)
}
