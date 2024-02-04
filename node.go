package flattenhtml

import (
	"errors"

	"golang.org/x/net/html"
)

// NodeIterator is a simple iterator that can iterate over a slice of *Node.
// It is used to iterate over the nodes that are flattened by a Flattener and
// perform different operations using the methods that are defined on the NodeIterator.
type NodeIterator struct {
	nodes       []*Node
	cursorIndex uint
}

// Node is a simple wrapper around *html.Node.
// It allows read/write operations on the *html.Node along with keeping the
// structure of the HTML tree.
type Node struct {
	htmlNode   *html.Node
	attributes map[string]string
	removed    bool
}

// FilterOption is a function that accepts a *Node and returns a boolean.
// The boolean value is true if the given *Node should be included in the
// NodeIterator and false otherwise.
type FilterOption func(node *Node) bool

var ErrParentlessNode = errors.New("node with no parent cannot be removed")

// NewNodeIterator creates a new NodeIterator.
func NewNodeIterator() *NodeIterator {
	return &NodeIterator{
		nodes: make([]*Node, 0),
	}
}

// Add adds the given *Node to the NodeIterator.
// This does not change the html.Node tree.
// It is expected that NodeIterator and Node are managed by the flattener.
func (n *NodeIterator) Add(node *Node) *NodeIterator {
	n.nodes = append(n.nodes, node)

	return n
}

// Len returns the number of nodes in the NodeIterator.
func (n *NodeIterator) Len() int {
	counter := 0

	for _, node := range n.nodes {
		if !node.IsRemoved() {
			counter++
		}
	}

	return counter
}

// Each iterates over the nodes in the NodeIterator and calls the given function.
func (n *NodeIterator) Each(f func(node *Node)) {
	for _, node := range n.nodes {
		if !node.IsRemoved() {
			f(node)
		}
	}
}

// First returns the first non-removed node in the NodeIterator.
// If there is no non-removed node, it returns nil.
func (n *NodeIterator) First() *Node {
	for _, node := range n.nodes {
		if !node.IsRemoved() {
			return node
		}
	}

	return nil
}

// Next iterates over the nodes in the NodeIterator and returns the next non-removed node.
// It starts from the first element of the NodeIterator and proceed to the next item on each
// call to Next. If there is no non-removed node, it returns nil.
// Once received nil, must be considered as the end of the iteration.
// Use Reset to start the iteration from the beginning.
func (n *NodeIterator) Next() *Node {
	for _, node := range n.nodes[n.cursorIndex:] {
		if !node.IsRemoved() {
			n.cursorIndex++

			return node
		}
	}

	return nil
}

// Reset resets the cursor index to the beginning of the NodeIterator.
func (n *NodeIterator) Reset() {
	n.cursorIndex = 0
}

// Filter filters the nodes in the NodeIterator using the given FilterOption.
// It returns a new NodeIterator that can iterate over the filtered nodes.
// For more complex filtering, you can use FilterOr or FilterAnd methods.
func (n *NodeIterator) Filter(option FilterOption) *NodeIterator {
	filteredNodes := NewNodeIterator()

	for _, node := range n.nodes {
		if node.IsRemoved() {
			continue
		}

		if option(node) {
			filteredNodes.Add(node)
		}
	}

	return filteredNodes
}

// FilterOr filters the nodes in the NodeIterator using the given FilterOptions.
// All the given options will be combined using OR operator. It means that if any
// of the given options returns true for a node, it will be included in the
// filtered NodeIterator and the rest of options will be ignored for that node.
func (n *NodeIterator) FilterOr(options ...FilterOption) *NodeIterator {
	filteredNodes := NewNodeIterator()

	for _, node := range n.nodes {
		if node.IsRemoved() {
			continue
		}

		for _, option := range options {
			if option(node) {
				filteredNodes.Add(node)

				break
			}
		}
	}

	return filteredNodes
}

// FilterAnd filters the nodes in the NodeIterator using the given FilterOptions.
// All the given options will be combined using AND operator. It means that if all
// the given options return true for a node, it will be included in the
// filtered NodeIterator. If any of the given options returns false for a node,
// the node will be filtered out and the rest of the options will be ignored for that node.
func (n *NodeIterator) FilterAnd(options ...FilterOption) *NodeIterator {
	filteredNodes := NewNodeIterator()

	for _, node := range n.nodes {
		if node.IsRemoved() {
			continue
		}

		filtered := true

		for _, option := range options {
			if !option(node) {
				filtered = false

				break
			}
		}

		if filtered {
			filteredNodes.Add(node)
		}
	}

	return filteredNodes
}

// NewNode creates a new Node with the given *html.Node.
func NewNode(htmlNode *html.Node) *Node {
	attrs := make(map[string]string, len(htmlNode.Attr))

	for _, attr := range htmlNode.Attr {
		attrs[attr.Key] = attr.Val
	}

	return &Node{
		htmlNode:   htmlNode,
		attributes: attrs,
	}
}

// IsRemoved returns true if the Node is removed from the NodeIterator
// and html.Node tree.
func (n *Node) IsRemoved() bool {
	return n.removed
}

// Remove removes the Node from the NodeIterator and html.Node tree.
// It won't be available if you use the NodeManager.Render.
func (n *Node) Remove() error {
	if n.htmlNode.Parent == nil {
		return ErrParentlessNode
	}

	n.htmlNode.Parent.RemoveChild(n.htmlNode)

	n.removed = true

	return nil
}

// TagName returns the tag name of the Node.
func (n *Node) TagName() string {
	return n.htmlNode.Data
}

// Attributes returns a map of strings containing attributes key and values of the Node.
func (n *Node) Attributes() map[string]string {
	return n.attributes
}

// Attribute returns the value of the given attribute key.
// The second return value is a boolean that indicates whether the given key is found.
func (n *Node) Attribute(key string) (string, bool) {
	val, ok := n.attributes[key]

	return val, ok
}

// SetAttribute sets the value of the given attribute key for the node.
// If the given key does not exist, it will be added to the node as a
// new attribute. Otherwise, the value of the given key will be updated.
func (n *Node) SetAttribute(key, value string) {
	exists := false

	for i, attr := range n.htmlNode.Attr {
		if attr.Key == key {
			n.htmlNode.Attr[i].Val = value

			exists = true

			break
		}
	}

	if !exists {
		n.htmlNode.Attr = append(n.htmlNode.Attr, html.Attribute{
			Key: key,
			Val: value,
		})
	}

	n.attributes[key] = value
}

// RemoveAttribute removes the given attribute key from the node.
// If the given key does not exist, it will be ignored.
func (n *Node) RemoveAttribute(key string) {
	for i, attr := range n.htmlNode.Attr {
		if attr.Key == key {
			n.htmlNode.Attr = append(n.htmlNode.Attr[:i], n.htmlNode.Attr[i+1:]...)

			break
		}
	}

	delete(n.attributes, key)
}

// HTMLNode returns the underlying *html.Node of the Node.
// Any write operation on the *html.Node might corrupt the structure of the HTML tree.
func (n *Node) HTMLNode() *html.Node {
	return n.htmlNode
}
