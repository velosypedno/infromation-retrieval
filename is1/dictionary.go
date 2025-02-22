package main

import "path/filepath"

func makeDictionary(args CommandArgs, wordSet map[string]struct{}) error {
	for _, directory := range args.directories {
		err := parseDirectory(directory, args, wordSet)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseDirectory(directory string, args CommandArgs, wordSet map[string]struct{}) error {
	files, err := getTxtFilePaths(directory, args)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		ext := filepath.Ext(filePath)
		if ext == ".fb2" {
			err := fb2Downloader(filePath, wordSet)
			if err != nil {
				return err
			}
		} else {
			err := defaultDownloader(filePath, wordSet)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
