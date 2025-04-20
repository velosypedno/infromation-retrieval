package producers

import (
	"is5/internal/core"
	"is5/internal/logger"
)

type IndexPairsProducer struct {
	IndexFilesToBePaired       chan string
	MergedIndexFilesToBePaired chan string
	IndexFilesPairs            chan core.Pair[string]
}

func (p *IndexPairsProducer) Produce() error {
	go func() {
		defer close(p.IndexFilesPairs)

		var buf []string

		for {
			select {
			case val, ok := <-p.IndexFilesToBePaired:
				logger.Log.Debug("Read from IndexFilesToBePaired cannel, doc: ", val)
				if !ok {
					p.IndexFilesToBePaired = nil
					logger.Log.Info("IndexFilesToBePaired cannel closed")
				}
				buf = append(buf, val)

			case val, ok := <-p.MergedIndexFilesToBePaired:
				logger.Log.Debug("Read from MergedIndexFilesToBePaired cannel, doc: ", val)
				if !ok {
					p.MergedIndexFilesToBePaired = nil
					logger.Log.Info("MergedIndexFilesToBePaired cannel closed")
				}
				buf = append(buf, val)
			default:
				if len(buf) < 2 {
					continue
				}
			}

			for len(buf) >= 2 {
				logger.Log.Debug("Send pair to IndexFilesPairs channel: ", buf[0], " ", buf[1], ", channel: ", p.IndexFilesPairs)
				p.IndexFilesPairs <- core.Pair[string]{
					First:  buf[0],
					Second: buf[1],
				}
				buf = buf[2:]
			}

			if p.IndexFilesToBePaired == nil && p.MergedIndexFilesToBePaired == nil && len(buf) < 2 {
				break
			}
		}
	}()
	return nil
}
