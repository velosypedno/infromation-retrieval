package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func saveDictionary(wordSet map[string]struct{}, args CommandArgs) error {
	extToLoader := map[string]func(wordSet map[string]struct{}, filePath string) error{
		".txt":  txtLoader,
		".json": jsonLoader,
		".csv":  csvLoader,
	}

	_, fileName := filepath.Split(args.outputPath)
	ext := filepath.Ext(fileName)

	loader, exists := extToLoader[ext]
	if !exists {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	err := loader(wordSet, args.outputPath)
	if err != nil {
		return err
	}
	return nil
}

func txtLoader(wordSet map[string]struct{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for word := range wordSet {
		_, err := writer.WriteString(word + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func jsonLoader(wordSet map[string]struct{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	words := make([]string, 0, len(wordSet))
	for word := range wordSet {
		words = append(words, word)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(words)
}

func csvLoader(wordSet map[string]struct{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	words := make([]string, 0, len(wordSet))
	for word := range wordSet {
		words = append(words, word)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(words); err != nil {
		return err
	}
	return nil

}

func defaultDownloader(filePath string, wordSet map[string]struct{}) error {
	regexPattern := `[a-zA-Zа-яА-ЯїЇєЄ]+([-'][a-zA-Zа-яА-ЯїЇєЄ]+)*`
	wordRegex := regexp.MustCompile(regexPattern)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := wordRegex.FindAllString(strings.ToLower(line), -1)
		for _, word := range tokens {
			wordSet[word] = struct{}{}
		}
	}
	return nil
}

func fb2Downloader(filePath string, wordSet map[string]struct{}) error {
	regexPattern := `[a-zA-Zа-яА-ЯїЇєЄ]+([-'][a-zA-Zа-яА-ЯїЇєЄ]+)*`
	wordRegex := regexp.MustCompile(regexPattern)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)

	var inBody bool
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "body" {
				inBody = true
			}
		case xml.EndElement:
			if t.Name.Local == "body" {
				inBody = false
			}
		case xml.CharData:
			if inBody {
				text := strings.ToLower(strings.TrimSpace(string(t)))
				words := wordRegex.FindAllString(text, -1)
				for _, word := range words {
					wordSet[word] = struct{}{}
				}
			}
		}
	}
	return nil
}
