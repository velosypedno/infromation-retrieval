package main

import (
	"log"
)

func main() {
	args := parseCommandLine()
	err := isValidArgs(args)
	if err != nil {
		log.Fatal(err)
	}
	wordSet := make(map[string]struct{})
	err = makeDictionary(args, wordSet)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Size of dictionary is %v words\n", len(wordSet))

	err = saveDictionary(wordSet, args)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Dictionary saved successfully to %v\n", args.outputPath)
}
