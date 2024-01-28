package flattenhtml

import (
	"context"
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

// Flattener is an interface for the logic that decides how the HTML tree
// should be traversed and flattened.
type Flattener interface {
	// Flatten is a callback function called for each node
	// in the HTML tree. It accepts a *html.Node as the argument and returns
	// an error if any. If the error is not nil, the iteration stops and the
	// error is returned.
	Flatten(node *html.Node) error

	// GetNodesByKey returns a NodeIterator that can iterate over the nodes
	// that are flattened using the flattener and filtered by the given key.
	// If the given key is not found in the flattened document, it returns
	// nil.
	GetNodesByKey(key string) *NodeIterator

	// IsMyType allows each flattener implementation to decide whether the given
	// Flattener is of the same type as itself or not.
	IsMyType(flattener Flattener) bool

	// Len the final number of categories or keys that were created by the flattener.
	Len() int
}

// NodeManager is an interface for the top-level logic of this package.
// This package is responsible to parse HTML nodes in some way, perform
// some modifications or read-only operations on them, and then render
// the HTML tree.
// There are different approaches to initiate a NodeManager:
//  1. NewNodeManagerFromReader: It accepts an io.Reader and parses
//     the HTML tree from it.
//  2. NewNodeManagerFromURL: It accepts a URL and parses the HTML
//     tree from the response body of the URL.
//  3. NewNodeManager: It accepts a *html.Node and uses it as the
//     root of the HTML tree.
//
// Using approaches 2 and 3 follow the [html.Parse] method to parse
// the HTML tree.
//
// [html.Parse]: https://pkg.go.dev/golang.org/x/net/html#Parse
type NodeManager struct {
	root *html.Node
}

// ErrNoFlattener is returned when no flattener is provided to the Parse method, or
// no flattener is found in the MultiCursor.
var ErrNoFlattener = errors.New("at least one flattener should be provided")

// NewNodeManager creates a new DefaultNodeManager with the given
// *html.Node as the root of the HTML tree.
func NewNodeManager(root *html.Node) *NodeManager {
	return &NodeManager{
		root: root,
	}
}

// NewNodeManagerFromReader creates a new DefaultNodeManager with
// the HTML tree parsed from the given io.Reader.
func NewNodeManagerFromReader(r io.Reader) (*NodeManager, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return NewNodeManager(root), nil
}

// NewNodeManagerFromURL creates a new DefaultNodeManager with the
// HTML tree parsed from the response body of the given URL.
func NewNodeManagerFromURL(ctx context.Context, url string) (*NodeManager, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return NewNodeManagerFromReader(resp.Body)
}

// Parse parses the HTML tree tha has been converted to *html.Node before.
// It accepts a set of Flattener that decides how the HTML tree should be
// traversed and flattened.
// If any of the flatteners returns an error, the iteration stops and the
// error is returned.
func (n *NodeManager) Parse(flatteners ...Flattener) (*MultiCursor, error) {
	if len(flatteners) == 0 {
		return nil, ErrNoFlattener
	}

	if err := nodeIterator(n.root, flatteners...); err != nil {
		return nil, err
	}

	return NewMultiCursor(flatteners...), nil
}

// Render renders the HTML tree to the given writer.
func (n *NodeManager) Render(w io.Writer) error {
	return html.Render(w, n.root)
}

// nodeIterator loops through all the *html.Node in the HTML tree.
// It continues until all the nodes are traversed.
// For each node that it meets, it calls the callback method of all
// the given flatteners to treat the node based on their logic.
func nodeIterator(node *html.Node, flatteners ...Flattener) error {
	if node == nil {
		return nil
	}

	for _, flattener := range flatteners {
		if err := flattener.Flatten(node); err != nil {
			return err
		}
	}

	if err := nodeIterator(node.FirstChild, flatteners...); err != nil {
		return err
	}

	return nodeIterator(node.NextSibling, flatteners...)
}
