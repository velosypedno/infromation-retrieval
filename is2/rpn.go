package main

import (
	"errors"
	"log"
	"math/big"
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
	log.Printf("query - %v", infixNotation)
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
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
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

func strToBigInt(str string) big.Int {
	bi, _ := new(big.Int).SetString(str, 2)
	return *bi
}

func executeQuery(query []string) (big.Int, error) {
	stack := []big.Int{}

	for len(query) > 0 {
		token := query[0]
		if o, ok := strToOperation[token]; ok {
			err := operationHandler(o, &stack)
			if err != nil {
				return big.Int{}, err
			}
		} else {
			stack = append(stack, strToBigInt(token))
		}
		query = query[1:]
	}
	if len(stack) > 1 || len(stack) == 0 {
		return big.Int{}, errors.New("wrong query syntaxes")
	}
	return stack[0], nil
}

func operationHandler(o Operation, stack *[]big.Int) error {
	switch o {
	case AND, OR:
		if len(*stack) < 2 {
			return errors.New("wrong query syntaxes")
		}
		firstOperand := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		secondOperand := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]

		var result big.Int
		if o == AND {
			result.And(&firstOperand, &secondOperand)
		} else {
			result.Or(&firstOperand, &secondOperand)
		}
		*stack = append(*stack, result)
	case NOT:
		if len(*stack) < 1 {
			return errors.New("wrong query syntaxes")
		}
		operand := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		*stack = append(*stack, *new(big.Int).Not(&operand))
	}
	return nil
}
