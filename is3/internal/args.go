package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	Dirs  []string
	Exts  []string
	Query string
}

func ParseArgs() (Args, error) {
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
	args.Dirs = strings.Split(*dirsFlag, ",")
	args.Exts = strings.Split(*extsFlag, ",")
	args.Query = *queryFlag
	err := validateArgs(args)
	return args, err

}

func validateArgs(args Args) error {
	if args.Query == "" {
		return errors.New("no query provided")
	}
	if len(args.Dirs) == 0 || (len(args.Dirs) == 1 && args.Dirs[0] == "") {
		return errors.New("no directories provided")
	}
	for _, path := range args.Dirs {
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
	if len(args.Exts) == 0 || (len(args.Exts) == 1 && args.Exts[0] == "") {
		return errors.New("no file extensions provided")
	}
	for _, ext := range args.Exts {
		if ext == "" || ext[0] != '.' {
			return fmt.Errorf("invalid file extension: %s", ext)
		}
	}

	return nil
}
