package args

import (
	"errors"
	"fmt"
	"os"
)

type ArgsValidator struct{}

func (v ArgsValidator) Validate(a Args) error {
	switch a.Mode {
	case Index:
		return validateIndexMode(a)
	case Search:
		return validateSearchMode(a)
	case IndexSearch:
		err := validateIndexMode(a)
		if err != nil {
			return err
		}
		err = validateSearchMode(a)
		return err
	}
	return nil
}

func validateSearchMode(a Args) error {
	if a.Query == "" {
		return errors.New("no query provided")
	}
	return nil
}

func validateIndexMode(a Args) error {
	if len(a.Dirs) == 0 || (len(a.Dirs) == 1 && a.Dirs[0] == "") {
		return errors.New("no directories provided")
	}
	for _, path := range a.Dirs {
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
	if len(a.Exts) == 0 || (len(a.Exts) == 1 && a.Exts[0] == "") {
		return errors.New("no file extensions provided")
	}
	for _, ext := range a.Exts {
		if ext == "" || ext[0] != '.' {
			return fmt.Errorf("invalid file extension: %s", ext)
		}
	}
	return nil
}
