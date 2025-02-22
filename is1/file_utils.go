package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getTxtFilePaths(directory string, args CommandArgs) ([]string, error) {
	minSize := args.minSize * 1024
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var textFiles []string
	extensions := make(map[string]bool)
	for _, ext := range args.extensions {
		extensions[ext] = true
	}
	for _, entry := range entries {
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if extensions[ext] {
			fileInfo, _ := entry.Info()
			if fileInfo.Size() < int64(minSize) {
				continue
			}
			log.Printf("File %v is suitable\n", filepath.Join(directory, entry.Name()))
			textFiles = append(textFiles, filepath.Join(directory, entry.Name()))
		}
	}
	return textFiles, nil
}
