package main

import (
	"is4/internal"
	"sort"
)

type PermIndex []struct {
	Perm string
	Term string
}

type PermIndexer struct {
	IndexWriter internal.Writer[PermIndex]
}

func (pi *PermIndexer) Index(terms []internal.TermDoc) error {
	permList := []permTerm{}
	for _, term := range terms {
		permutations := generatePermutations(term.Term + "$")
		for _, permutation := range permutations {
			permList = append(permList, permTerm{perm: permutation, term: term.Term})
		}
	}
	permMap := make(map[string]string)

	for _, permTerm := range permList {
		if _, ok := permMap[permTerm.perm]; !ok {
			permMap[permTerm.perm] = permTerm.term
		}
	}
	var permIndex PermIndex
	for perm, term := range permMap {
		permIndex = append(permIndex, struct {
			Perm string
			Term string
		}{Perm: perm, Term: term})
	}
	sort.Slice(permIndex, func(i, j int) bool {
		return permIndex[i].Perm < permIndex[j].Perm
	})

	err := pi.IndexWriter.Write(&permIndex)
	if err != nil {
		return err
	}
	return nil
}

type permTerm struct {
	perm string
	term string
}

func generatePermutations(term string) []string {
	permutations := []string{}
	for i := 0; i < len(term); i++ {
		permutations = append(permutations, term[i:]+term[:i])
	}
	return permutations
}
