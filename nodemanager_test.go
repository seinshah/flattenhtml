package flattenhtml_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestNewNodeManager(t *testing.T) {
	t.Parallel()

	sampleHTML := "<html><head></head><body></body></html>"

	testCases := []struct {
		name    string
		result  func() (*flattenhtml.NodeManager, error)
		wantErr bool
	}{
		{
			name: "new node manager from html.Node",
			result: func() (*flattenhtml.NodeManager, error) {
				sampleTree, err := html.Parse(strings.NewReader(sampleHTML))
				if err != nil {
					return nil, err
				}

				return flattenhtml.NewNodeManager(sampleTree), nil
			},
		},
		{
			name: "new node manager from reader",
			result: func() (*flattenhtml.NodeManager, error) {
				return flattenhtml.NewNodeManagerFromReader(strings.NewReader(sampleHTML))
			},
		},
		{
			name: "new node manager from url",
			result: func() (*flattenhtml.NodeManager, error) {
				return flattenhtml.NewNodeManagerFromURL(context.Background(), "https://google.com")
			},
		},
		{
			name: "non-existing url",
			result: func() (*flattenhtml.NodeManager, error) {
				return flattenhtml.NewNodeManagerFromURL(context.Background(), "https://test.non-existing")
			},
			wantErr: true,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := tc.result()

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

func TestNodeManager_Parse(t *testing.T) {
	t.Parallel()

	sampleHTML := "<html><head></head><body></body></html>"

	testCases := []struct {
		name      string
		flattener []flattenhtml.Flattener
		// expectedCallbacks is the expected number of times the callback
		// method of each flattener is called. exclude html tag from the count.
		expectedCallbacks []int
		wantErr           bool
		errType           error
	}{
		{
			name:      "parse with no flattener",
			flattener: nil,
			wantErr:   true,
			errType:   flattenhtml.ErrNoFlattener,
		},
		{
			name:      "parse with one flattener that trigger error",
			flattener: []flattenhtml.Flattener{&sampleFlattener{withErr: true}},
			wantErr:   true,
			errType:   errSample,
		},
		{
			name:              "parse with one flattener",
			flattener:         []flattenhtml.Flattener{&sampleFlattener{}},
			expectedCallbacks: []int{4},
		},
		{
			name:      "parse with two flatteners and one of them trigger error",
			flattener: []flattenhtml.Flattener{&sampleFlattener{}, &sampleFlattener{withErr: true}},
			wantErr:   true,
			errType:   errSample,
		},
		{
			name:              "parse with two flatteners",
			flattener:         []flattenhtml.Flattener{&sampleFlattener{}, &sampleFlattener{}},
			expectedCallbacks: []int{4, 4},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := flattenhtml.NewNodeManagerFromReader(strings.NewReader(sampleHTML))
			require.NoError(t, err)

			mc, err := resp.Parse(tc.flattener...)

			if tc.wantErr {
				require.Error(t, err)

				if tc.errType != nil {
					require.IsType(t, tc.errType, err)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, mc)

			for index, f := range tc.flattener {
				require.Equal(t, tc.expectedCallbacks[index], f.Len())
			}
		})
	}
}

func TestNodeManager_Render(t *testing.T) {
	t.Parallel()

	sampleHTML := "<html><head></head><body></body></html>"

	testCases := []struct {
		name     string
		expected string
		wantErr  bool
		errType  error
	}{
		{
			name:     "successful render",
			expected: sampleHTML,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := flattenhtml.NewNodeManagerFromReader(strings.NewReader(sampleHTML))
			require.NoError(t, err)

			mc, err := resp.Parse(&sampleFlattener{})

			require.NoError(t, err)
			require.NotNil(t, mc)

			rendered := bytes.Buffer{}

			err = resp.Render(&rendered)

			if tc.wantErr {
				require.Error(t, err)

				if tc.errType != nil {
					require.IsType(t, tc.errType, err)
				}

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, rendered.String())
		})
	}
}
