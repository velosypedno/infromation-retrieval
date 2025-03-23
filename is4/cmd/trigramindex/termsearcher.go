package main

import (
	"is4/internal"
	"strings"
)

type TermSearcher struct {
	IndexReader internal.Reader[TrigramIndex]
	PostFilter  internal.Filter[string, string]
}

func (t TermSearcher) Search(term string) ([]string, error) {
	template := term
	if term == "" {
		return nil, nil
	}
	index, err := t.IndexReader.Read()
	if err != nil {
		return nil, err
	}
	if term[0] != '*' {
		term = "$" + term
	}
	if term[len(term)-1] != '*' {
		term = term + "$"
	}

	tokens := strings.FieldsFunc(term, func(r rune) bool {
		return r == '*'
	})

	trigrams := []string{}
	for _, token := range tokens {
		trigrams = append(trigrams, generateTrigrams(token)...)
	}
	if len(trigrams) == 0 {
		return nil, nil
	}

	terms := termsByTrigramBinarySearch(trigrams[0], index)
	for i, trigram := range trigrams {
		if i == 0 {
			continue
		}
		newTerms := termsByTrigramBinarySearch(trigram, index)
		terms = *termsIntersection(&terms, &newTerms)
	}
	terms = *t.PostFilter.Filter(&terms, template)
	return terms, nil
}

func termsByTrigramBinarySearch(trigram string, index *TrigramIndex) []string {
	start := 0
	end := len(*index) - 1
	for start <= end {
		mid := (start + end) / 2
		if trigram < (*index)[mid].Trigram {
			end = mid - 1
		} else if trigram > (*index)[mid].Trigram {
			start = mid + 1
		} else {
			return (*index)[mid].Terms
		}
	}
	return nil
}

func termsIntersection(terms1 *[]string, terms2 *[]string) *[]string {
	t1 := 0
	t2 := 0
	result := []string{}
	for t1 < len(*terms1) && t2 < len(*terms2) {
		if (*terms1)[t1] < (*terms2)[t2] {
			t1++
		} else if (*terms1)[t1] > (*terms2)[t2] {
			t2++
		} else {
			result = append(result, (*terms1)[t1])
			t1++
			t2++
		}
	}
	return &result
}
