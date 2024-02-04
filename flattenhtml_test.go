package flattenhtml_test

import (
	"fmt"

	"github.com/seinshah/flattenhtml"
	"golang.org/x/net/html"
)

type sampleFlattener struct {
	called      int
	withErr     bool
	defaultKeys []string
}

var _ flattenhtml.Flattener = (*sampleFlattener)(nil)

var errSample = fmt.Errorf("sample error")

func (s *sampleFlattener) Flatten(_ *html.Node) error {
	if s.withErr {
		return errSample
	}

	s.called++

	return nil
}

func (s *sampleFlattener) GetNodesByKey(key string) *flattenhtml.NodeIterator {
	for _, k := range s.defaultKeys {
		if k == key {
			return flattenhtml.NewNodeIterator()
		}
	}

	return nil
}

func (s *sampleFlattener) IsMyType(_ flattenhtml.Flattener) bool {
	return true
}

func (s *sampleFlattener) Len() int {
	return s.called
}
