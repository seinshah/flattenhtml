package flattenhtml_test

import (
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestTagFlattener(t *testing.T) {
	t.Parallel()

	sampleNodes := []*html.Node{
		{
			Type: html.ElementNode,
			Data: "div",
		},
		{
			Type: html.DocumentNode,
			Data: "a",
		},
	}

	flattener := flattenhtml.NewTagFlattener()

	for _, node := range sampleNodes {
		err := flattener.Flatten(node)

		require.NoError(t, err)
	}

	require.Equal(t, 1, flattener.Len())
	require.Equal(t, 1, flattener.GetNodesByKey("div").Len())
	require.Nil(t, flattener.GetNodesByKey("a"))
	require.False(t, flattener.IsMyType(&sampleFlattener{}))
}
