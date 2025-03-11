package main

import (
	"strconv"
	"unicode"
)

func isOperator(token string) bool {
	if len(token) < 2 {
		return false
	}
	if token[0] == '/' {
		for _, char := range token[1:] {
			if !unicode.IsDigit(char) {
				return false
			}
		}
		return true
	}
	return false
}

func infixToRPN(infixNotation []string) []string {
	var output []string
	currentOperator := ""
	for _, token := range infixNotation {
		if isOperator(token) {
			if currentOperator != "" {
				output = append(output, currentOperator)
			}
			currentOperator = token
		} else {
			output = append(output, token)
		}
	}
	if currentOperator != "" {
		output = append(output, currentOperator)
	}
	return output
}

func ExecuteQuery(query []string, posIndex *PositionalIndex) ([]int, error) {
	query = infixToRPN(query)
	stack := make([]TermIndex, 0)
	for _, token := range query {
		if isOperator(token) {
			term2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			term1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			distance, err := strconv.Atoi(token[1:])
			if err != nil {
				return nil, err
			}
			stack = append(stack, positionalIntersect(term1, term2, distance))
		} else {
			index := binarySearch(token, posIndex)
			if index == -1 {
				return nil, nil
			}
			stack = append(stack, posIndex.InvIndex[index])
		}
	}
	docs := make([]int, 0)
	for _, doc := range stack[0].Docs {
		docs = append(docs, doc.Doc)
	}
	return docs, nil
}

func positionalIntersect(term1 TermIndex, term2 TermIndex, distance int) TermIndex {
	intersectTerm := TermIndex{Term: term1.Term + " " + term2.Term, Docs: make([]struct {
		Doc  int
		Poss []int
	}, 0)}
	docP1 := 0
	docP2 := 0
	for docP1 < len(term1.Docs) && docP2 < len(term2.Docs) {
		if term1.Docs[docP1].Doc == term2.Docs[docP2].Doc {
			posP1 := 0
			posP2 := 0
			created := false
			positions1 := term1.Docs[docP1].Poss
			positions2 := term2.Docs[docP2].Poss

			for posP1 < len(positions1) && posP2 < len(positions2) {
				if positions1[posP1] > positions2[posP2] {
					posP2++
					continue
				}
				if positions2[posP2]-positions1[posP1] <= distance {
					if !created {
						intersectTerm.Docs = append(intersectTerm.Docs, struct {
							Doc  int
							Poss []int
						}{Doc: term1.Docs[docP1].Doc, Poss: make([]int, 0)})
						created = true
					}
					intersectTerm.Docs[len(intersectTerm.Docs)-1].Poss =
						append(intersectTerm.Docs[len(intersectTerm.Docs)-1].Poss, positions2[posP2])
					posP2++
				} else {
					posP1++
				}
			}
			created = false
			docP1++
			docP2++
		} else if term1.Docs[docP1].Doc < term2.Docs[docP2].Doc {
			docP1++
		} else {
			docP2++
		}
	}
	return intersectTerm
}
