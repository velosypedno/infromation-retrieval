package internal

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

type GobWriter[T any] struct {
	FileName string
}

func (w *GobWriter[T]) Write(data *T) error {
	dir := filepath.Dir(w.FileName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(w.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(*data)
	if err != nil {
		return err
	}
	return nil
}
