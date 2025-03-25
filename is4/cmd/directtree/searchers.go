package main

import (
	"errors"
	"is4/internal"
	"strings"

	"github.com/google/btree"
)

type StringItem struct {
	Term string
}

func (a StringItem) Less(b btree.Item) bool {
	return a.Term < b.(StringItem).Term
}

func (ti *TreeIndex) ToBTree() *btree.BTree {
	tree := btree.New(2)
	for _, term := range ti.Terms {
		tree.ReplaceOrInsert(StringItem{Term: term})
	}
	return tree
}

type TreeSearcher struct {
	IndexReader internal.Reader[TreeIndex]
}

func (t TreeSearcher) Search(template string) ([]string, error) {
	result := []string{}
	if template == "" {
		return nil, nil
	}

	isValid := true
	for _, char := range template[0 : len(template)-1] {
		if char == '*' {
			isValid = false
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid query")
	}

	treeIndex, err := t.IndexReader.Read()
	if err != nil {
		return nil, err
	}
	bTree := treeIndex.ToBTree()

	isPrefixSearch := template[len(template)-1] == '*'
	if isPrefixSearch {
		bTree.Ascend(func(i btree.Item) bool {
			prefix := template[:len(template)-1]
			term := i.(StringItem).Term
			if strings.HasPrefix(term, prefix) {
				result = append(result, term)
			}
			return true
		})
	} else {
		bTree.Ascend(func(i btree.Item) bool {
			term := i.(StringItem).Term
			if term == template {
				result = append(result, term)
			}
			return true
		})
	}
	return result, nil

}
