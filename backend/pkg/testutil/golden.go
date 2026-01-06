package testutil

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func AssertGoldenJSON(t *testing.T, golden string, got []byte) {
	t.Helper()

	if got == nil {
		assert.Equal(t, "", golden)
		return
	}

	// JSONの整形
	var gotJSON interface{}
	err := json.Unmarshal(got, &gotJSON)
	require.NoError(t, err)
	got, err = json.MarshalIndent(gotJSON, "", "  ")
	require.NoError(t, err)

	if *update {
		err := os.MkdirAll(filepath.Dir(golden), 0755)
		require.NoError(t, err)
		err = os.WriteFile(golden, got, 0644)
		require.NoError(t, err)
	}

	want, err := os.ReadFile(golden)
	require.NoError(t, err)

	require.JSONEq(t, string(want), string(got))
}
