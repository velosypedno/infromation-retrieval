package main

import (
	"is3/internal"
	"log"
	"regexp"
)

func main() {
	args, err := internal.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	indexFiles(args)

	re := regexp.MustCompile(`'[^']*'|\S+`)

	query := re.FindAllString(args.Query, -1)
	posIndex := loadIndex()

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
