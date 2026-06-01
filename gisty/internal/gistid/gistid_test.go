package gistid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	require.Equal(t, "ABC abc 123", Sanitize("ABC abc 123\n<>"))
}
