package main

import (
	"encoding/csv"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"
)

func indexOf(slice []int, value int) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func buildMatrix(termsSet map[TermPosition]struct{}, indexToPath map[int]string) error {
	terms := make(map[string]struct{})
	files := make(map[int]struct{})

	for tp := range termsSet {
		terms[tp.term] = struct{}{}
		files[tp.fileIndex] = struct{}{}
	}

	termList := make([]string, 0, len(terms))
	for t := range terms {
		termList = append(termList, t)
	}
	sort.Strings(termList)

	fileIndexes := make([]int, 0, len(files))
	for f := range files {
		fileIndexes = append(fileIndexes, f)
	}
	sort.Ints(fileIndexes)
	fileNames := make([]string, 0, len(files))
	for _, fileIndex := range fileIndexes {
		fileNames = append(fileNames, indexToPath[fileIndex])
	}

	matrix := make([][]string, len(termList)+1)
	matrix[0] = append([]string{"Term"}, fileNames...)

	termIndex := make(map[string]int)
	for i, term := range termList {
		termIndex[term] = i + 1
		matrix[i+1] = make([]string, len(fileIndexes)+1)
		matrix[i+1][0] = term
	}

	for i := 1; i < len(termList)+1; i++ {
		for j := 1; j < len(fileIndexes)+1; j++ {
			matrix[i][j] = "0"
		}
	}

	for tp := range termsSet {
		row := termIndex[tp.term]
		col := indexOf(fileIndexes, tp.fileIndex) + 1
		matrix[row][col] = "1"
	}

	file, err := os.Create("matrix.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, row := range matrix {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	log.Println("Matrix saved")
	return nil
}

func loadCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err == nil {
		log.Printf("File - %v, size - %v", file.Name(), fileInfo.Size())
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func parseVector(matrix [][]string, index int) string {
	return strings.Join(matrix[index][1:], "")
}

func binarySearch(matrix [][]string, target string) int {
	low, high := 1, len(matrix)-1

	for low <= high {
		mid := (low + high) / 2
		word := matrix[mid][0]

		if word == target {
			return mid
		} else if word < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func leftPad(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(padChar), length-len(s)) + s
}

func getFilesByVector(vector *big.Int, matrix [][]string) []string {
	vectorLen := len(matrix[0]) - 1
	bitMask := make([]bool, 0, vectorLen)

	vectorStr := (*vector).Text(2)
	vectorStr = leftPad(vectorStr, vectorLen, '0')
	for _, char := range vectorStr {
		if char == '0' {
			bitMask = append(bitMask, false)
		} else {
			bitMask = append(bitMask, true)
		}
	}
	fileNames := []string{}
	for index, filename := range matrix[0][1:] {
		if bitMask[index] {
			fileNames = append(fileNames, filename)
		}
	}
	return fileNames
}

func replaceQueryTermsWithVector(query *[]string, matrix [][]string) {
	for index, token := range *query {
		switch token {
		case "AND", "OR", "NOT":
		default:
			rowIndex := binarySearch(matrix, token)
			if rowIndex == -1 {
				(*query)[index] = "0"
			} else {
				(*query)[index] = parseVector(matrix, rowIndex)
			}
		}
	}
}

func searchByMatrix(query string) {
	start := time.Now()
	matrix, err := loadCSV("./matrix.csv")
	if err != nil {
		log.Fatal(err)
	}

	queryRPN := infixToRPN(strings.Split(query, " "))
	replaceQueryTermsWithVector(&queryRPN, matrix)
	result, err := executeQuery(queryRPN)
	if err != nil {
		log.Fatal(err)
	}
	fileNames := getFilesByVector(&result, matrix)
	log.Printf("Files - %v\n", fileNames)

	elapsed := time.Since(start)
	log.Printf("Execution duration: %d milliseconds\n", elapsed.Milliseconds())
}
