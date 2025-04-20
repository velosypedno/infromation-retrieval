package producers

import (
	"is5/internal/core"
	"is5/internal/logger"
)

type DocsFromTxtToChannelProducer struct {
	FilePath   string
	DocChannel chan core.Doc
	DocReader  core.Reader[[]core.Doc]
}

func (p *DocsFromTxtToChannelProducer) Produce() error {
	docs, err := p.DocReader.Read(p.FilePath)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			close(p.DocChannel)
			logger.Log.Debug("DocsFromTxtToChannelProducer closed: ", p.DocChannel)
		}()
		for _, doc := range docs {
			logger.Log.Debug("Send doc to channel: ", doc.Path, " ", doc.Id)
			p.DocChannel <- doc
		}
	}()

	return nil
}
