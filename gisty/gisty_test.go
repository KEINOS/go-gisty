package gisty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeGistID(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		input    string
		expected string
	}{
		{input: "7101f542be23e5048198e2a27c3cfda8", expected: "7101f542be23e5048198e2a27c3cfda8"},
		{input: "abc<>def", expected: "abcdef"},
		{input: "abcあいうえおde💩f", expected: "abcdef"},
		{input: "abc\tde\nf", expected: "abcdef"},
	} {
		expect := test.expected
		actual := SanitizeGistID(test.input)

		require.Equal(t, expect, actual,
			"test #%d: input: %s", index+1, test.input)
	}
}
