package flattenhtml_test

import (
	"testing"

	"github.com/seinshah/flattenhtml"
	"github.com/stretchr/testify/require"
)

func TestMultiCursor_SelectCursor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		flatteners []flattenhtml.Flattener
		flattener  flattenhtml.Flattener
		wantErr    bool
		errType    error
	}{
		{
			name:       "select no flattener",
			flatteners: nil,
			flattener:  nil,
			wantErr:    true,
			errType:    flattenhtml.ErrNoFlattener,
		},
		{
			name:       "select non-existing flattener",
			flatteners: nil,
			flattener:  &sampleFlattener{},
			wantErr:    true,
			errType:    flattenhtml.ErrNoFlattener,
		},
		{
			name:       "select successfully",
			flatteners: []flattenhtml.Flattener{&sampleFlattener{}},
			flattener:  &sampleFlattener{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mc := flattenhtml.NewMultiCursor(tc.flatteners...)

			cu, err := mc.SelectCursor(tc.flattener)

			if tc.wantErr {
				require.Error(t, err)

				if tc.errType != nil {
					require.IsType(t, tc.errType, err)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, cu)
		})
	}
}

func TestCursor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		flattener   flattenhtml.Flattener
		key         string
		exists      bool
		expectedLen int
	}{
		{
			name:        "select absent key",
			flattener:   &sampleFlattener{called: 5},
			key:         "absent",
			expectedLen: 5,
		},
		{
			name:        "select existing key",
			flattener:   &sampleFlattener{defaultKeys: []string{"existing"}},
			key:         "existing",
			exists:      true,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cu := flattenhtml.NewMultiCursor(tc.flattener).First()

			require.NotNil(t, cu)
			require.NotNil(t, cu.SelectNodes(tc.key))
		})
	}
}
