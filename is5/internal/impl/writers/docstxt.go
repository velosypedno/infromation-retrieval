package writers

import (
	"is5/internal/core"
	"is5/internal/logger"
	"os"
	"strconv"
)

type DocsTxtWriter struct {
	FilePath string
}

func (w DocsTxtWriter) Write(docs []core.Doc) error {
	logger.Log.Info("Writing docs to ", w.FilePath)
	file, err := os.OpenFile(w.FilePath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		_, err := file.WriteString(doc.Path + " " + strconv.Itoa(doc.Id) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
