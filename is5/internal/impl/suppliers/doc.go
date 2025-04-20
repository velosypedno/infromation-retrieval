package suppliers

import "is5/internal/core"

type DocSupplier struct {
	DocCannel chan core.Doc
}

func (s *DocSupplier) Supply() (core.Doc, error) {
	doc, ok := <-s.DocCannel
	if !ok {
		return core.Doc{}, nil
	}
	return doc, nil

}
