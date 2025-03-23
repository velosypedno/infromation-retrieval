package internal

import (
	"encoding/gob"
	"os"
)

type GobReader[T any] struct {
	FileName string
}

func (r *GobReader[T]) Read() (*T, error) {
	f, err := os.Open(r.FileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	var data T
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
