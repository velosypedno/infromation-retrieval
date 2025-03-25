package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	dirs  []string
	exts  []string
	query string
}

func parseArgs() Args {
	dirsFlag := flag.String("dirs", "", "Directories to scan (comma-separated, e.g. /path1,/path2)")
	extsFlag := flag.String("exts", ".txt", "Allowed file extensions (comma-separated, e.g. .txt,.md)")
	queryFlag := flag.String("query", "", "Query like 'help OR NOT give'")

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go [OPTIONS]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println(`  go run main.go -dirs="/data/texts" -exts=".txt,.md" -query="NOT apple"`)
	}

	var args Args
	flag.Parse()
	args.dirs = strings.Split(*dirsFlag, ",")
	args.exts = strings.Split(*extsFlag, ",")
	args.query = *queryFlag
	return args

}

func validateArgs(args Args) error {
	if args.query == "" {
		return errors.New("no query provided")
	}
	if len(args.dirs) == 0 || (len(args.dirs) == 1 && args.dirs[0] == "") {
		return errors.New("no directories provided")
	}
	for _, path := range args.dirs {
		if path == "" {
			return errors.New("empty directory path found")
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			msg := fmt.Sprintf("%s is not directory\n", path)
			return errors.New(msg)
		}
	}
	if len(args.exts) == 0 || (len(args.exts) == 1 && args.exts[0] == "") {
		return errors.New("no file extensions provided")
	}
	for _, ext := range args.exts {
		if ext == "" || ext[0] != '.' {
			return fmt.Errorf("invalid file extension: %s", ext)
		}
	}

	return nil
}
