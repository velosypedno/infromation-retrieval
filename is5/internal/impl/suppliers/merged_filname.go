package suppliers

import (
	"is5/internal/logger"
	"strconv"
	"sync/atomic"
)

type MergedIndexFilenameSupplier struct {
	BaseDir      string
	StartIndexId int64
}

func (i *MergedIndexFilenameSupplier) Supply() (string, error) {
	id := atomic.AddInt64(&i.StartIndexId, 1)
	logger.Log.Info("Supply merged index filename: ", i.BaseDir+"/"+strconv.FormatInt(id, 10)+".merged")
	return i.BaseDir + "/" + strconv.FormatInt(id, 10) + ".merged", nil
}
