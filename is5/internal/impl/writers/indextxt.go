package writers

import (
	"is5/internal/core"
	"os"
	"strconv"
	"strings"
)

type IndexTxtWriter struct {
	FilePathSupplier core.Supplier[string]
	FileNamesCh      chan string
}

func (w IndexTxtWriter) Write(index *[]core.TermToDocIds) error {
	filePath, err := w.FilePathSupplier.Supply()
	defer func() { w.FileNamesCh <- filePath }()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, termToDocIds := range *index {

		entry := termToDocIds.Term + " " + intSliceToString(termToDocIds.DocIds) + "\n"
		_, err := file.WriteString(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func intSliceToString(ids []int) string {
	var builder strings.Builder
	for i, num := range ids {
		if i > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.Itoa(num))
	}
	return builder.String()
}
