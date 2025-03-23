package main

import (
	"is4/internal"
	"sort"
)

type TrigramIndex []struct {
	Trigram string
	Terms   []string
}

type TrigramIndexer struct {
	IndexWriter internal.Writer[TrigramIndex]
}

func (ti *TrigramIndexer) Index(terms []internal.TermDoc) error {
	trigramList := []trigramTerm{}
	for _, term := range terms {
		extendedWord := "$" + term.Term + "$"
		trigrams := generateTrigrams(extendedWord)
		for _, trigram := range trigrams {
			trigramList = append(trigramList, trigramTerm{trigram: trigram, terms: term.Term})
		}
	}
	trigramMap := make(map[string]map[string]struct{})

	for _, trigramTerm := range trigramList {
		if _, ok := trigramMap[trigramTerm.trigram]; !ok {
			trigramMap[trigramTerm.trigram] = make(map[string]struct{})
		}
		trigramMap[trigramTerm.trigram][trigramTerm.terms] = struct{}{}
	}
	var trigrams TrigramIndex
	for trigram, termsMap := range trigramMap {
		terms := make([]string, 0, len(termsMap))
		for term := range termsMap {
			terms = append(terms, term)
		}
		trigrams = append(trigrams, struct {
			Trigram string
			Terms   []string
		}{Trigram: trigram, Terms: terms})
	}
	for _, trigram := range trigrams {
		sort.Slice(trigram.Terms, func(i, j int) bool {
			return trigram.Terms[i] < trigram.Terms[j]
		})
	}
	sort.Slice(trigrams, func(i, j int) bool {
		return trigrams[i].Trigram < trigrams[j].Trigram
	})
	err := ti.IndexWriter.Write(&trigrams)
	if err != nil {
		return err
	}
	return nil
}

type trigramTerm struct {
	trigram string
	terms   string
}

func generateTrigrams(term string) []string {
	trigrams := []string{}

	for i := 0; i < len(term)-2; i++ {
		trigrams = append(trigrams, term[i:i+3])
	}
	return trigrams
}
