package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CommandArgs struct {
	directories []string
	extensions  []string
	minSize     int
	outputPath  string
}

func parseCommandLine() CommandArgs {
	dirsFlag := flag.String("dirs", "", "Directories to scan (comma-separated, e.g. /path1,/path2)")
	extsFlag := flag.String("exts", ".txt", "Allowed file extensions (comma-separated, e.g. .txt,.md)")
	minSizeFlag := flag.Int("minsize", 0, "Min file size in kb")
	outFlag := flag.String("out", "./dictionary.txt", "Output file path")

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go [OPTIONS]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println(`  go run main.go -dirs="/data/texts" -exts=".txt,.md" -minsize=102400 -out="words.txt"`)
	}

	var args CommandArgs
	flag.Parse()
	args.directories = strings.Split(*dirsFlag, ",")
	args.extensions = strings.Split(*extsFlag, ",")
	args.minSize = *minSizeFlag
	args.outputPath = *outFlag
	return args

}

func isValidArgs(args CommandArgs) error {
	if len(args.directories) == 0 || (len(args.directories) == 1 && args.directories[0] == "") {
		return errors.New("no directories provided")
	}
	for _, path := range args.directories {
		if path == "" {
			return errors.New("empty directory path found")
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			message := fmt.Sprintf("%s is not directory\n", path)
			return errors.New(message)
		}
	}
	if len(args.extensions) == 0 || (len(args.extensions) == 1 && args.extensions[0] == "") {
		return errors.New("no file extensions provided")
	}
	for _, ext := range args.extensions {
		if ext == "" || ext[0] != '.' {
			return fmt.Errorf("invalid file extension: %s", ext)
		}
	}
	if args.minSize < 0 {
		return errors.New("file size cannot be negative")
	}
	outDir, fileName := filepath.Split(args.outputPath)
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("output directory cannot be created: %s", outDir)
	}
	file, err := os.CreateTemp(outDir, fileName)
	if err != nil {
		return fmt.Errorf("%s cannot be created in directory %s", fileName, outDir)
	}
	file.Close()
	os.Remove(file.Name())

	return nil
}
