package flattenhtml_test

import (
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestWithTag(t *testing.T) {
	t.Parallel()

	sampleNode := flattenhtml.NewNode(&html.Node{Data: "div"})

	testCases := []struct {
		name     string
		tag      string
		expected bool
	}{
		{
			name:     "non-existing tag",
			tag:      "a",
			expected: false,
		},
		{
			name:     "existing tag",
			tag:      "div",
			expected: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := flattenhtml.WithTag(tc.tag)(sampleNode)

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestWithAttribute(t *testing.T) {
	t.Parallel()

	sampleNode := flattenhtml.NewNode(&html.Node{
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "test",
			},
		},
	})

	testCases := []struct {
		name     string
		attrKey  string
		expected bool
	}{
		{
			name:     "non-existing attribute",
			attrKey:  "id",
			expected: false,
		},
		{
			name:     "existing attribute",
			attrKey:  "class",
			expected: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := flattenhtml.WithAttribute(tc.attrKey)(sampleNode)

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestWithAttributeValueAs(t *testing.T) {
	t.Parallel()

	sampleNode := flattenhtml.NewNode(&html.Node{
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "test",
			},
		},
	})

	testCases := []struct {
		name     string
		attrKey  string
		attrVal  string
		expected bool
	}{
		{
			name:     "non-existing attribute",
			attrKey:  "id",
			attrVal:  "test",
			expected: false,
		},
		{
			name:     "existing attribute with different value",
			attrKey:  "class",
			attrVal:  "test2",
			expected: false,
		},
		{
			name:     "existing attribute with same value",
			attrKey:  "class",
			attrVal:  "test",
			expected: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := flattenhtml.WithAttributeValueAs(tc.attrKey, tc.attrVal)(sampleNode)

			require.Equal(t, tc.expected, actual)
		})
	}
}
