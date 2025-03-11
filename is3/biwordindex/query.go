package main

import (
	"strings"
)

type Operation int

const (
	AND Operation = iota
	OR
	NOT
)

var operationStrings = []string{"AND", "OR", "NOT"}
var strToOperation = map[string]Operation{
	AND.String(): AND,
	OR.String():  OR,
	NOT.String(): NOT,
}

func (o Operation) String() string {
	if int(o) < len(operationStrings) {
		return operationStrings[o]
	}
	return "UNKNOWN"
}

var precedence = map[string]int{
	NOT.String(): 3,
	AND.String(): 2,
	OR.String():  1,
}

func infixToRPN(infixNotation []string) []string {
	var output []string
	var stack []string

	for _, token := range infixNotation {
		switch token {
		case AND.String(), NOT.String(), OR.String():
			for len(stack) > 0 && precedence[token] <= precedence[stack[len(stack)-1]] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		default:
			output = append(output, token)
		}
	}
	for len(stack) != 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output
}

func ExecuteQuery(query []string, posIndex *PositionalIndex) ([]int, error) {
	query = infixToRPN(query)
	stack := make([]TermIndex, 0)
	for _, token := range query {
		switch token {
		case AND.String():
			term2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			term1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, and(term1, term2))
		case OR.String():
			term2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			term1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, or(term1, term2))
		case NOT.String():
			stack = append(stack, not(stack[len(stack)-1], posIndex))
		default:
			if token[0] != '\'' {
				index := binarySearch(token, posIndex)
				if index == -1 {
					return nil, nil
				}
				stack = append(stack, posIndex.InvIndex[index])
			} else {
				phrase := token[1 : len(token)-1]
				terms := strings.Split(phrase, " ")
				biStack := make([]TermIndex, 0)
				for i, token := range terms {
					if i == 0 {
						continue
					}
					token = terms[i-1] + " " + token
					index := binarySearch(token, posIndex)
					if index == -1 {
						return nil, nil
					}
					biStack = append(biStack, posIndex.InvIndex[index])
				}
				for len(biStack) > 1 {
					term2 := biStack[len(biStack)-1]
					biStack = biStack[:len(biStack)-1]
					term1 := biStack[len(biStack)-1]
					biStack = biStack[:len(biStack)-1]
					biStack = append(biStack, and(term1, term2))
				}
				stack = append(stack, biStack[0])
			}
		}
	}
	return stack[0].Docs, nil
}

func and(term1 TermIndex, term2 TermIndex) TermIndex {
	intersectTerm := TermIndex{Term: term1.Term + " " + term2.Term, Docs: make([]int, 0)}
	docP1 := 0
	docP2 := 0
	for docP1 < len(term1.Docs) && docP2 < len(term2.Docs) {
		if term1.Docs[docP1] == term2.Docs[docP2] {
			intersectTerm.Docs = append(intersectTerm.Docs, term1.Docs[docP1])
			docP1++
			docP2++
		} else if term1.Docs[docP1] < term2.Docs[docP2] {
			docP1++
		} else {
			docP2++
		}
	}
	return intersectTerm
}

func or(term1 TermIndex, term2 TermIndex) TermIndex {
	intersectTerm := TermIndex{Term: term1.Term + " " + term2.Term, Docs: make([]int, 0)}
	docP1 := 0
	docP2 := 0
	for docP1 < len(term1.Docs) && docP2 < len(term2.Docs) {
		if term1.Docs[docP1] <= term2.Docs[docP2] {
			intersectTerm.Docs = append(intersectTerm.Docs, term1.Docs[docP1])
			docP1++
		} else {
			intersectTerm.Docs = append(intersectTerm.Docs, term2.Docs[docP2])
			docP2++
		}
	}
	for docP1 < len(term1.Docs) {
		intersectTerm.Docs = append(intersectTerm.Docs, term1.Docs[docP1])
		docP1++
	}
	for docP2 < len(term2.Docs) {
		intersectTerm.Docs = append(intersectTerm.Docs, term2.Docs[docP2])
		docP2++
	}

	return intersectTerm
}

func not(term TermIndex, posIndex *PositionalIndex) TermIndex {
	notDocs := make(map[int]struct{})
	for _, doc := range term.Docs {
		notDocs[doc] = struct{}{}
	}
	intersectTerm := TermIndex{Term: "!" + term.Term, Docs: make([]int, 0)}
	for doc := range posIndex.Docs {
		if _, ok := notDocs[doc]; ok {
			continue
		}
		intersectTerm.Docs = append(intersectTerm.Docs, doc)
	}
	return intersectTerm
}
