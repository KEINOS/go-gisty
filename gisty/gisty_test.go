package gisty

import (
	"testing"

	"github.com/stretchr/testify/require"
	_ "golang.org/x/oauth2"
)

const (
	sanitizeTestGistID   = "7101f542be23e5048198e2a27c3cfda8"
	sanitizeExpectedText = "abcdef"
)

func TestSanitizeGistID(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		input    string
		expected string
	}{
		{input: sanitizeTestGistID, expected: sanitizeTestGistID},
		{input: "abc<>def", expected: sanitizeExpectedText},
		{input: "abcあいうえおde💩f", expected: sanitizeExpectedText},
		{input: "abc\tde\nf", expected: sanitizeExpectedText},
	} {
		expect := test.expected
		actual := SanitizeGistID(test.input)

		require.Equal(t, expect, actual,
			"test #%d: input: %s", index+1, test.input)
	}
}
