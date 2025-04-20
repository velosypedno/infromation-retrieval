package args

import (
	"flag"
	"fmt"
	"is5/internal/core"
	"strings"
)

type Mode int

const (
	Index Mode = iota
	Search
	IndexSearch
)

type Args struct {
	Dirs  []string
	Exts  []string
	Query string
	Mode  Mode
}

type ArgsSupplier struct {
	Validator core.Validator[Args]
}

func (s *ArgsSupplier) Supply() (Args, error) {
	dirsFlag := flag.String("dirs", "", "Directories to scan (comma-separated, e.g. /path1,/path2)")
	extsFlag := flag.String("exts", "", "Allowed file extensions (comma-separated, e.g. .txt,.md)")
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
	if (args.Query != "") && (args.Exts[0] == "") && (args.Dirs[0] == "") {
		args.Mode = Search
	} else if len(args.Query) > 0 && ((len(args.Exts) > 0) || (len(args.Dirs) > 0)) {
		args.Mode = IndexSearch
	} else {
		args.Mode = Index
	}

	err := s.Validator.Validate(args)
	return args, err

}
