package docreaders

import (
	"bufio"
	"is5/internal/logger"
	"log"
	"os"
)

type LineByLineReader struct {
	ChanBuffSize int
}

func (d LineByLineReader) Read(path string) (<-chan string, error) {
	logger.Log.Info("Reading file line by line: ", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out := make(chan string, d.ChanBuffSize)
	go func() {
		defer file.Close()
		defer close(out)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return out, nil

}
