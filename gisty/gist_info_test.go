package gisty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGistInfo_empty_input(t *testing.T) {
	t.Parallel()

	gistInfo, err := NewGistInfo("")

	require.Error(t, err, "empty line should return error")
	require.Empty(t, gistInfo, "it should return an empty GistInfo on error")
	require.Contains(t, err.Error(), "empty line", "error should contain the reason")
}

func TestNewGistInfo_invalid_input(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		name   string
		input  string
		reason string // msg that should be contained in the error
	}{
		{
			name:   "files field is not in 'n files' format",
			input:  "d5b9800c	Title of gist item 2	6 filetes	secret	2022-04-16T06:08:46Z",
			reason: "failed to parse number of files from:",
		},
		{
			name:   "updatedAt field is not in RFC3339 format",
			input:  "d5b9800c	Title of gist item 2	6 files	secret	2022-09-18 18:56:10 +0000 UTC",
			reason: "failed to parse time from:",
		},
		{
			name:   "updateAt field is missing",
			input:  "d5b9800c	Title of gist item 2	6 files	secret",
			reason: "missing number of chunks:",
		},
		{
			name:   "isPublic field is not 'public' or 'secret'",
			input:  "d5b9800c	Title of gist item 2	6 files	private	2022-04-16T06:08:46Z	foo	bar",
			reason: "failed to parse isPublic from:",
		},
	} {
		gistInfo, err := NewGistInfo(test.input)

		require.Error(t, err,
			"test #%d: \"%s\" failed\ninvalid format should be an error", index+1, test.name)
		require.Empty(t, gistInfo,
			"test #%d: \"%s\" failed\nit should return an empty GistInfo on error", index+1, test.name)
		require.Contains(t, err.Error(), test.reason,
			"test #%d: \"%s\" failed\nerror should contain the reason", index+1, test.name)
	}
}
