package main

import (
	"is3/internal"
	"log"
	"strings"
)

func main() {
	args, err := internal.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	indexFiles(args)
	posIndex := loadIndex()
	query := strings.Split(args.Query, " ")
	queryResults, err := ExecuteQuery(query, &posIndex)
	if err != nil {
		log.Fatal(err)
	}
	docs := make([]string, 0)
	for _, doc := range queryResults {
		docs = append(docs, posIndex.Docs[doc])
	}
	log.Println(docs)
}
