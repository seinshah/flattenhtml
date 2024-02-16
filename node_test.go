package flattenhtml_test

import (
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestNodeIterator(t *testing.T) {
	t.Parallel()

	nodes := []*flattenhtml.Node{
		flattenhtml.NewNode(&html.Node{
			Data: "div",
			Type: html.ElementNode,
		}),
		flattenhtml.NewNode(&html.Node{
			Data: "a",
			Type: html.ElementNode,
		}),
	}

	nodeIterator := flattenhtml.NewNodeIterator()

	for _, node := range nodes {
		nodeIterator.Add(node)
	}

	require.Len(t, nodes, nodeIterator.Len())

	activeIndex := 0

	nodeIterator.Each(func(node *flattenhtml.Node) {
		require.Equal(t, nodes[activeIndex].TagName(), node.TagName())

		activeIndex++
	})

	filteredIterator := nodeIterator.Filter(func(node *flattenhtml.Node) bool {
		return node.TagName() == "div"
	})

	require.Equal(t, 1, filteredIterator.Len())
	require.Equal(t, "div", filteredIterator.First().TagName())

	orFilteredIterator := nodeIterator.FilterOr(
		func(node *flattenhtml.Node) bool {
			return node.TagName() == "div"
		},
		func(node *flattenhtml.Node) bool {
			return node.TagName() == "a"
		},
	)

	require.Equal(t, 2, orFilteredIterator.Len())
	require.Equal(t, "div", orFilteredIterator.Next().TagName())
	require.Equal(t, "a", orFilteredIterator.Next().TagName())

	orFilteredIterator.Reset()
	require.Equal(t, "div", orFilteredIterator.Next().TagName())

	andFilteredIterator := nodeIterator.FilterAnd(
		func(node *flattenhtml.Node) bool {
			return node.TagName() == "div"
		},
		func(node *flattenhtml.Node) bool {
			return node.TagName() == "a"
		},
	)

	require.Equal(t, 0, andFilteredIterator.Len())
}

func TestNode(t *testing.T) {
	t.Parallel()

	node := flattenhtml.NewNode(&html.Node{
		Data: "div",
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "test",
			},
		},
	})

	require.Len(t, node.Attributes(), 1, "expected number of attributes")

	attrVal, ok := node.Attribute("class")
	require.True(t, ok, "class attribute exists")
	require.Equal(t, "test", attrVal, "class attribute has expected value")

	attrVal, ok = node.Attribute("non-existent")
	require.False(t, ok, "non-existent attribute doesn't exist")
	require.Equal(t, "", attrVal, "non-existent attribute has empty value")

	// Setting a new attribute and test the same process again
	node.SetAttribute("new-attr", node.TagName())
	attrVal, ok = node.Attribute("new-attr")
	require.True(t, ok, "new-attr attribute exists")
	require.Equal(t, node.TagName(), attrVal, "new-attr attribute has expected value")
	require.Len(t, node.Attributes(), 2, "expected number of attributes after new-attr")

	// Removing the new attribute and test the same process again
	node.RemoveAttribute("new-attr")
	attrVal, ok = node.Attribute("new-attr")
	require.False(t, ok, "new-attr attribute doesn't exist")
	require.Equal(t, "", attrVal, "new-attr attribute has empty value")
	require.Len(t, node.Attributes(), 1, "expected number of attributes after removing new-attr")

	hn := node.HTMLNode()
	require.NotNil(t, hn, "expected *html.Node")
	require.Equal(t, node.TagName(), hn.Data, "expected tag name")
}

func TestNode_Remove(t *testing.T) {
	t.Parallel()

	// Creating the representation of <body><div></div><a></a></body>
	bodyNode := &html.Node{
		Data: "body",
		Type: html.ElementNode,
	}

	divNode := &html.Node{
		Data:   "div",
		Type:   html.ElementNode,
		Parent: bodyNode,
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "test-div",
			},
		},
	}

	aNode := &html.Node{
		Data:        "a",
		Type:        html.ElementNode,
		Parent:      divNode,
		PrevSibling: divNode,
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "test-a",
			},
		},
	}

	divNode.NextSibling = aNode

	nodeIterator := flattenhtml.NewNodeIterator()

	nodeIterator.Add(flattenhtml.NewNode(bodyNode))
	nodeIterator.Add(flattenhtml.NewNode(divNode))
	nodeIterator.Add(flattenhtml.NewNode(aNode))

	nodeIterator.Each(func(node *flattenhtml.Node) {
		if node.TagName() == aNode.Data {
			err := node.Remove()

			require.NoError(t, err)
		}
	})

	require.Equal(t, 2, nodeIterator.Len())

	nodeIterator.Each(func(node *flattenhtml.Node) {
		require.NotEqual(t, aNode.Data, node.TagName())
	})

	fnTagName := func(node *flattenhtml.Node) bool {
		return node.TagName() == aNode.Data
	}
	fnAttrClass := func(node *flattenhtml.Node) bool {
		attr, _ := node.Attribute("class")

		return attr == "test-a"
	}

	require.Equal(t, 0, nodeIterator.Filter(fnTagName).Len())
	require.Equal(t, 0, nodeIterator.FilterOr(fnTagName, fnAttrClass).Len())
	require.Equal(t, 0, nodeIterator.FilterAnd(fnTagName, fnAttrClass).Len())
}

func TestNode_AddNewNode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		operation func(node *flattenhtml.Node) *flattenhtml.Node
		validate  func(baseNode, newNode *flattenhtml.Node) bool
	}{
		{
			name: "append child to a node",
			operation: func(node *flattenhtml.Node) *flattenhtml.Node {
				return node.AppendChild(
					flattenhtml.NodeTypeElement,
					"span",
					nil,
				)
			},
			validate: func(baseNode, newNode *flattenhtml.Node) bool {
				return baseNode.HTMLNode().LastChild == newNode.HTMLNode() &&
					newNode.HTMLNode().Parent == baseNode.HTMLNode()
			},
		},
		{
			name: "prepend child to a node",
			operation: func(node *flattenhtml.Node) *flattenhtml.Node {
				return node.PrependChild(
					flattenhtml.NodeTypeElement,
					"span",
					nil,
				)
			},
			validate: func(baseNode, newNode *flattenhtml.Node) bool {
				return baseNode.HTMLNode().FirstChild == newNode.HTMLNode() &&
					newNode.HTMLNode().Parent == baseNode.HTMLNode()
			},
		},
		{
			name: "append sibling to a node",
			operation: func(node *flattenhtml.Node) *flattenhtml.Node {
				return node.AppendSibling(
					flattenhtml.NodeTypeElement,
					"span",
					nil,
				)
			},
			validate: func(baseNode, newNode *flattenhtml.Node) bool {
				return baseNode.HTMLNode().NextSibling == newNode.HTMLNode() &&
					newNode.HTMLNode().PrevSibling == baseNode.HTMLNode() &&
					baseNode.HTMLNode().Parent == newNode.HTMLNode().Parent
			},
		},
		{
			name: "prepend sibling to a node",
			operation: func(node *flattenhtml.Node) *flattenhtml.Node {
				return node.PrependSibling(
					flattenhtml.NodeTypeElement,
					"span",
					nil,
				)
			},
			validate: func(baseNode, newNode *flattenhtml.Node) bool {
				return baseNode.HTMLNode().PrevSibling == newNode.HTMLNode() &&
					newNode.HTMLNode().NextSibling == baseNode.HTMLNode() &&
					baseNode.HTMLNode().Parent == newNode.HTMLNode().Parent
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			div := &html.Node{
				Data: "div",
				Type: html.ElementNode,
			}

			a := &html.Node{
				Data: "a",
				Type: html.ElementNode,
			}

			label := &html.Node{
				Data: "label",
				Type: html.ElementNode,
			}

			div.AppendChild(a)
			a.AppendChild(label)

			baseNode := flattenhtml.NewNode(a)

			newNode := tc.operation(baseNode)

			require.True(t, tc.validate(baseNode, newNode))
		})
	}
}
