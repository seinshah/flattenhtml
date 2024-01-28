package flattenhtml

// WithTag is a function that filters Node based on their tag name.
// If the node's tag name is the same is the given tag, it will be included in
// the final output.
func WithTag(tag string) FilterOption {
	return func(node *Node) bool {
		return node.TagName() == tag
	}
}

// WithAttribute returns a FilterOption that filters nodes by the given key.
// The Node will be included in the final output if it has an attribute with the given key.
func WithAttribute(key string) FilterOption {
	return func(node *Node) bool {
		_, ok := node.Attribute(key)

		return ok
	}
}

// WithAttributeValueAs returns a FilterOption that filters nodes by the given key and value.
// The Node will be included in the final output if it has an attribute with the given key
// and the value of that attribute is equal to the given value.
func WithAttributeValueAs(key, value string) FilterOption {
	return func(node *Node) bool {
		val, ok := node.Attribute(key)

		if !ok || val != value {
			return false
		}

		return true
	}
}
