package mergers

import (
	"errors"
	"is5/internal/core"
	"is5/internal/logger"
	"log"
	"os"
	"strconv"
	"strings"
)

type DefaultIndexMerger struct {
	LineByLineReader      core.DocReader
	IndexFileNameSupplier core.Supplier[string]
	FileRemover           core.Remover[string]
}

func (d DefaultIndexMerger) Merge(filepath1 string, filepath2 string) (string, error) {
	filename, err := d.IndexFileNameSupplier.Supply()
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	lines1, err := d.LineByLineReader.Read(filepath1)
	if err != nil {
		return "", err
	}
	lines2, err := d.LineByLineReader.Read(filepath2)
	if err != nil {
		return "", err
	}

	var line1, line2 string
	var ok1, ok2 bool
	term1, term2 := "", ""
	var validTerm1, validTerm2 bool

	logger.Log.Debug("Merging files: ", filepath1, " and ", filepath2)
	for {
		if term1 == "" {
			line1, ok1 = <-lines1
			if !ok1 {
				break
			}
			term1, _, validTerm1 = strings.Cut(line1, " ")
			if !validTerm1 {
				logger.Log.Error("Invalid line format, line: ", line1, ", file: ", filepath1)
				return "", errors.New("invalid line format")
			}
		}
		if term2 == "" {
			line2, ok2 = <-lines2
			if !ok2 {
				break
			}
			term2, _, validTerm2 = strings.Cut(line2, " ")
			if !validTerm2 {
				logger.Log.Error("Invalid line format, line: ", line2, ", file: ", filepath2)
				return "", errors.New("invalid line format")
			}
		}

		if term1 < term2 {
			_, err = file.WriteString(line1 + "\n")
			if err != nil {
				return "", err
			}
			term1 = ""
		} else if term1 > term2 {
			_, err = file.WriteString(line2 + "\n")
			if err != nil {
				return "", err
			}
			term2 = ""
		} else {
			is1, err := docIdsFromString(line1)
			if err != nil {
				return "", err
			}
			is2, err := docIdsFromString(line2)
			if err != nil {
				return "", err
			}
			ids := mergeInts(is1, is2)
			entry := term1 + " " + intSliceToString(ids) + "\n"
			_, err = file.WriteString(entry)
			if err != nil {
				return "", err
			}
			term1 = ""
			term2 = ""
		}
	}
	for ok1 {
		_, err = file.WriteString(line1 + "\n")
		if err != nil {
			return "", err
		}
		line1, ok1 = <-lines1
	}
	for ok2 {
		_, err = file.WriteString(line2 + "\n")
		if err != nil {
			return "", err
		}
		line2, ok2 = <-lines2
	}
	log.Println("finished merging")
	err = d.FileRemover.Remove(filepath1)
	if err != nil {
		return "", err
	}
	err = d.FileRemover.Remove(filepath2)
	if err != nil {
		return "", err
	}
	return filename, nil

}

func mergeInts(ids1 []int, ids2 []int) []int {
	i, j := 0, 0
	result := make([]int, 0, len(ids1)+len(ids2))

	for i < len(ids1) && j < len(ids2) {
		if ids1[i] < ids2[j] {
			result = append(result, ids1[i])
			i++
		} else if ids1[i] > ids2[j] {
			result = append(result, ids2[j])
			j++
		} else {
			result = append(result, ids1[i])
			i++
			j++
		}
	}

	if i < len(ids1) {
		result = append(result, ids1[i:]...)
	}
	if j < len(ids2) {
		result = append(result, ids2[j:]...)
	}

	return result
}

func docIdsFromString(line string) ([]int, error) {
	_, idsString, valid := strings.Cut(line, " ")
	if !valid {
		return nil, errors.New("invalid line format")
	}
	idsStringSlice := strings.Split(idsString, " ")
	return stringSliceToIntSlice(idsStringSlice), nil
}

func stringSliceToIntSlice(stringSlice []string) []int {
	intSlice := make([]int, len(stringSlice))
	for i, str := range stringSlice {
		intSlice[i], _ = strconv.Atoi(str)
	}
	return intSlice
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
