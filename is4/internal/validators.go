package internal

import (
	"errors"
	"fmt"
	"os"
)

type ArgsValidator struct{}

func (v ArgsValidator) Validate(args Args) error {
	switch args.Mode {
	case Index:
		return validateIndexMode(args)
	case Search:
		return validateSearchMode(args)
	case IndexSearch:
		err := validateIndexMode(args)
		if err != nil {
			return err
		}
		err = validateSearchMode(args)
		return err
	}
	return nil
}

func validateSearchMode(args Args) error {
	if args.Query == "" {
		return errors.New("no query provided")
	}
	return nil
}

func validateIndexMode(args Args) error {
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
