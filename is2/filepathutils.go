package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getFilesFromDirectory(dir string, extsMap map[string]struct{}) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if _, ok := extsMap[ext]; ok {
			log.Printf("File %v is suitable\n", filepath.Join(dir, entry.Name()))
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}

func getFilesFromDirectories(dirs []string, exts []string) ([]string, error) {
	files := []string{}
	extsMap := make(map[string]struct{})
	for _, ext := range exts {
		extsMap[ext] = struct{}{}
	}

	for _, dir := range dirs {
		filesFromDir, err := getFilesFromDirectory(dir, extsMap)
		if err != nil {
			return nil, err
		}
		files = append(files, filesFromDir...)
	}
	return files, nil

}
