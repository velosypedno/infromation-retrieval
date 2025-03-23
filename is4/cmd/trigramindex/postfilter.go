package main

import (
	"strings"
)

type ByTemplateFiler struct{}

func (f ByTemplateFiler) Filter(terms *[]string, key string) *[]string {
	filteredTerms := []string{}
	justWord := true
	for _, char := range key {
		if char == '*' {
			justWord = false
			break
		}
	}
	if justWord {
		for _, term := range *terms {
			if term == key {
				filteredTerms = append(filteredTerms, term)
			}
		}
		return &filteredTerms
	}
	for _, term := range *terms {
		if checkByTemplate(term, key) {
			filteredTerms = append(filteredTerms, term)
		}
	}
	return &filteredTerms
}

func checkByTemplate(term string, template string) bool {
	if template[0] == '*' {
		return strings.HasSuffix(term, template[1:])
	}
	if template[len(template)-1] == '*' {
		return strings.HasPrefix(term, template[:len(template)-1])
	}
	sufAndPref := strings.FieldsFunc(template, func(r rune) bool {
		return r == '*'
	})
	if len(sufAndPref) != 2 {
		return false
	}
	isCorrectPrefix := strings.HasPrefix(term, sufAndPref[0])
	isCorrectSuffix := strings.HasSuffix(term, sufAndPref[1])
	return isCorrectPrefix && isCorrectSuffix
}
