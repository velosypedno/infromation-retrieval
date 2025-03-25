package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

type InvertedIndex map[string][]int

func buildIndex(termsSet map[TermPosition]struct{}, indexToPath map[int]string) error {
	invertedIndex := make(InvertedIndex)

	for termPos := range termsSet {
		if _, ok := invertedIndex[termPos.term]; !ok {
			invertedIndex[termPos.term] = []int{}
		}
		invertedIndex[termPos.term] = append(invertedIndex[termPos.term], termPos.fileIndex)
	}

	file, err := os.Create("inverted_index.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(invertedIndex); err != nil {
		return err
	}

	file, err = os.Create("index_to_path.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder = json.NewEncoder(file)
	if err := encoder.Encode(indexToPath); err != nil {
		return err
	}

	return nil
}

func searchByIndex(query string) {
	start := time.Now()
	index, err := loadInvertedIndex()
	if err != nil {
		return
	}
	queryRPN := infixToRPN(strings.Split(query, " "))
	replaceQueryTermsWithVectorFromIndex(&queryRPN, index)
	vector, err := executeQuery(queryRPN)
	if err != nil {
		log.Fatal(err)
	}
	files := getFilesByVectorFromIndex(&vector)
	log.Printf("files - %v", files)

	elapsed := time.Since(start)
	log.Printf("Execution duration: %d milliseconds\n", elapsed.Milliseconds())
}

func maxIndex() (int, error) {
	maxKey := -1
	filePath := "./index_to_path.json"

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return maxKey, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return maxKey, err
	}
	for key := range data {
		if num, err := strconv.Atoi(key); err == nil {
			if num > maxKey {
				maxKey = num
			}
		}
	}
	return maxKey, nil

}

func loadInvertedIndex() (InvertedIndex, error) {
	index := make(InvertedIndex)

	filePath := "./inverted_index.json"
	file, _ := os.Open(filePath)
	defer file.Close()
	info, _ := file.Stat()
	log.Printf("File - %v, size - %v", info.Name(), info.Size())
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var tempIndex InvertedIndex
	err = json.Unmarshal(fileData, &tempIndex)
	if err != nil {
		return nil, err
	}
	for key, value := range tempIndex {
		index[key] = value
	}

	return index, err
}

func replaceQueryTermsWithVectorFromIndex(query *[]string, index InvertedIndex) {
	vectorLength, err := maxIndex()
	if err != nil {
		return
	}
	for i, token := range *query {
		switch token {
		case "AND", "OR", "NOT":
		default:
			vector := make([]string, vectorLength+1)
			for j := range vector {
				vector[j] = "0"
			}
			if indexes, ok := index[token]; ok {
				for _, j := range indexes {
					vector[j] = "1"
				}
			}
			(*query)[i] = strings.Join(vector, "")
		}
	}
}

func loadJSONToMap() (map[string]string, error) {
	filePath := "./index_to_path.json"
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(fileData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getFilesByVectorFromIndex(vector *big.Int) []string {
	vectorLen, err := maxIndex()
	if err != nil {
		return []string{}
	}
	vectorLen++
	vectorStr := (*vector).Text(2)
	vectorStr = leftPad(vectorStr, vectorLen, '0')
	fileNames := []string{}
	indexToPath, err := loadJSONToMap()
	if err != nil {
		return []string{}
	}
	for index, char := range vectorStr {
		if char == '1' {
			fileNames = append(fileNames, indexToPath[fmt.Sprint(index)])
		}
	}
	return fileNames

}
