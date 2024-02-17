package flattenhtml

// MultiCursor is a helper struct that holds all the configured flatteners.
// It will usually be initiated by the NodeManager using the configured
// flatteners which can be later filtered to a single flattener using
// *MultiCursor.SelectFlattener method.
type MultiCursor struct {
	flatteners []Flattener
}

// Cursor is a helper struct that holds the selected flattener from the MultiCursor.
// It allows the caller to perform different operations on the flattened document using
// the selected flattener by *MultiCursor.SelectFlattener method.
type Cursor struct {
	flattener Flattener
}

// NewMultiCursor returns a new MultiCursor initiated by the NodeManager.
// This holds all the configured flatteners that are used separately to
// flatten the HTML tree.
// To perform the variety of operations on the flattened documents, first you need
// to select your desired flattener cursor using methods defined on MultiCursor.
func NewMultiCursor(flatteners ...Flattener) *MultiCursor {
	return &MultiCursor{
		flatteners: flatteners,
	}
}

// First returns the first Cursor from the MultiCursor initiated by the NodeManager.
// This Cursor will hold the reference to the first flattener you configured for
// the NodeManager.Parse method.
// If MultiCursor has no cursor, the result will be nil.
func (m *MultiCursor) First() *Cursor {
	if len(m.flatteners) == 0 {
		return nil
	}

	return &Cursor{flattener: m.flatteners[0]}
}

// SelectCursor returns a new Cursor with the selected flattener from the MultiCursor
// initiated by the NodeManager.
// If the given flattener is not found in the MultiCursor, it returns ErrNoFlattener.
func (m *MultiCursor) SelectCursor(flattener Flattener) (*Cursor, error) {
	if flattener == nil {
		return nil, ErrNoFlattener
	}

	var newFlattener Flattener

	for _, f := range m.flatteners {
		if f.IsMyType(flattener) {
			newFlattener = f

			break
		}
	}

	if newFlattener == nil {
		return nil, ErrNoFlattener
	}

	return &Cursor{flattener: newFlattener}, nil
}

// RegisterNewNode is used to add a newly and manually added nodes by the user to the cycle.
// It calls flatten method of all it's flatteners by giving the Node's underlying html.Node.
// New node can only be accessed by the NodeIterator and Cursor, if it is added to the cycle
// using this method.
func (m *MultiCursor) RegisterNewNode(node *Node) error {
	if len(m.flatteners) == 0 {
		return ErrNoFlattener
	}

	for _, f := range m.flatteners {
		if err := f.Flatten(node.htmlNode); err != nil {
			return err
		}
	}

	return nil
}

// SelectNodes returns a new NodeIterator that can iterates over the nodes that are selected
// by the given key and perform different operations.
// If the given key is not found in the flattened document, nodeIterator will have a zero length.
func (c *Cursor) SelectNodes(key string) *NodeIterator {
	nodes := c.flattener.GetNodesByKey(key)

	if nodes == nil {
		nodes = NewNodeIterator()
	}

	return nodes
}

// Len returns the final number of categories or keys that were created by the flattener.
func (c *Cursor) Len() int {
	return c.flattener.Len()
}

// RegisterNewNode is used to add a newly and manually added nodes by the user to the cycle.
// It calls flatten method of the cursor's flatteners by giving the Node's underlying html.Node.
// New node can only be accessed by the NodeIterator and Cursor, if it is added to the cycle
// using this method.
func (c *Cursor) RegisterNewNode(node *Node) error {
	return c.flattener.Flatten(node.HTMLNode())
}
