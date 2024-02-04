package flattenhtml_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestIntegration(t *testing.T) {
	rawHTML := `<html><body><div><p class="p1">hello</p><p class="p2">world</p></div></body></html>`
	expectedRawHTML := `<html><head></head><body><div><p class="p1">hello</p></div></body></html>`

	manager, err := flattenhtml.NewNodeManagerFromReader(strings.NewReader(rawHTML))
	require.NoError(t, err)
	require.NotNil(t, manager)

	mc, err := manager.Parse(flattenhtml.NewTagFlattener())
	require.NoError(t, err)
	require.NotNil(t, mc)

	cursor := mc.First()
	require.NotNil(t, cursor)

	// There should be 5 tags: html, head, body, div, p
	require.Equal(t, 5, cursor.Len())
	require.Equal(t, 0, cursor.SelectNodes("a").Len())

	pNodes := cursor.SelectNodes("p")
	require.Equal(t, 2, pNodes.Len())

	targetP := pNodes.Filter(flattenhtml.WithAttributeValueAs("class", "p2"))
	require.Equal(t, 1, targetP.Len())

	targetP.Each(func(node *flattenhtml.Node) {
		err = node.Remove()
		require.NoError(t, err)
	})

	require.Equal(t, 1, pNodes.Len())

	output := bytes.Buffer{}

	err = manager.Render(&output)
	require.NoError(t, err)

	require.Equal(t, expectedRawHTML, output.String())
}

func ExampleTagFlattener() {
	rawHTML := `<html><body><div><p class="p1">hello</p><p class="p2">world</p></div></body></html>`

	manager, err := flattenhtml.NewNodeManagerFromReader(strings.NewReader(rawHTML))
	if err != nil {
		panic(err)
	}

	mc, err := manager.Parse(flattenhtml.NewTagFlattener())
	if err != nil {
		panic(err)
	}

	cursor := mc.First()
	targetP := cursor.SelectNodes("p").Filter(flattenhtml.WithAttributeValueAs("class", "p2"))

	targetP.Each(func(node *flattenhtml.Node) {
		err = node.Remove()
		if err != nil {
			panic(err)
		}
	})

	output := bytes.Buffer{}

	err = manager.Render(&output)
	if err != nil {
		panic(err)
	}

	fmt.Println(output.String())

	// Output: <html><head></head><body><div><p class="p1">hello</p></div></body></html>
}
