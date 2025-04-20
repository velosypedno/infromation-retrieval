package suppliers

import (
	"log"
	"strconv"
	"sync/atomic"
)

type IndexFilenameSupplier struct {
	BaseDir      string
	StartIndexId int64
}

func (i *IndexFilenameSupplier) Supply() (string, error) {
	id := atomic.AddInt64(&i.StartIndexId, 1)
	if id%10 == 0 {
		log.Println("Index of file with id was created: ", id)
	}
	return i.BaseDir + "/" + strconv.FormatInt(id, 10) + ".txt", nil
}
