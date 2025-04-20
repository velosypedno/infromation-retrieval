package readers

import (
	"bufio"
	"errors"
	"is5/internal/core"
	"is5/internal/logger"
	"os"
	"strconv"
	"strings"
)

type TxtDocsReader struct{}

func (r TxtDocsReader) Read(path string) ([]core.Doc, error) {
	logger.Log.Info("Reading docs from ", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	docs := make([]core.Doc, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			return nil, errors.New("invalid line format")
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		docs = append(docs, core.Doc{
			Path: parts[0],
			Id:   id,
		})
	}
	return docs, nil
}
