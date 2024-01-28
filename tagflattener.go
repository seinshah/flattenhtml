package flattenhtml

import (
	"golang.org/x/net/html"
)

// TagFlattener is a Flattener that flattens the HTML tree by the tag name.
// When the NodeManager is initialized with this flattener, it will categorize
// NodeIterator by the tag name. Therefore, you can access all nodes with the
// same tag name (i.e., meta, a, p, etc.) using the GetNodesByKey method or
// Cursor.SelectNodes method.
type TagFlattener struct {
	flattened map[string]*NodeIterator
}

var _ Flattener = (*TagFlattener)(nil)

// NewTagFlattener creates a new TagFlattener.
func NewTagFlattener() *TagFlattener {
	return &TagFlattener{
		flattened: make(map[string]*NodeIterator),
	}
}

// Flatten is a callback function called for each node during the
// NodeManager.Parse. It will continue to categorize all nodes in their tag
// NodeIterator as NodeManager traverses the HTML tree. This method does not
// return an error.
func (t *TagFlattener) Flatten(node *html.Node) error {
	if node.Type == html.ElementNode {
		if _, ok := t.flattened[node.Data]; !ok {
			t.flattened[node.Data] = NewNodeIterator()
		}

		t.flattened[node.Data].Add(NewNode(node))
	}

	return nil
}

func (t *TagFlattener) GetNodesByKey(key string) *NodeIterator {
	return t.flattened[key]
}

func (t *TagFlattener) IsMyType(flattener Flattener) bool {
	_, ok := flattener.(*TagFlattener)

	return ok
}

// Len for tagflattener gives you the concrete number of tags in the HTML tree.
func (t *TagFlattener) Len() int {
	return len(t.flattened)
}
