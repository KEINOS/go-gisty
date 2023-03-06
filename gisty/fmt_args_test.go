package gisty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fmtArgs(t *testing.T) {
	t.Parallel()

	require.Empty(t, fmtArgs(),
		"no args should return empty string")

	require.Equal(t, "10 lines found", fmtArgs("%v lines found", 10),
		"if the first arg is a string, it should be formatted with the rest of the args")

	require.Equal(t, "1 2 3", fmtArgs(1, 2, 3),
		"if the first arg is not a string, it should return the concatenation of all args")
}
