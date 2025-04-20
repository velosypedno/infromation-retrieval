package removers

import "os"

type FileRemover struct{}

func (f FileRemover) Remove(filename string) error {
	return os.Remove(filename)
}
