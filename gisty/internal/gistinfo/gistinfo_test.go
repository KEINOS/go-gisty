package gistinfo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	info, err := Parse("abc\tdescription\t2 files\tpublic\t2026-06-01T00:00:00Z")

	require.NoError(t, err)
	require.Equal(t, "abc", info.GistID)
	require.Equal(t, "description", info.Description)
	require.Equal(t, 2, info.Files)
	require.True(t, info.IsPublic)
	require.Equal(t, "2026-06-01 00:00:00 +0000 UTC", info.UpdatedAt.String())
}

func TestParse_error(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "missing field", input: "abc"},
		{name: "invalid file count", input: "abc\tdescription\tmany files\tpublic\t2026-06-01T00:00:00Z"},
		{name: "invalid visibility", input: "abc\tdescription\t2 files\tprivate\t2026-06-01T00:00:00Z"},
		{name: "invalid timestamp", input: "abc\tdescription\t2 files\tsecret\tnow"},
	} {
		_, err := Parse(test.input)
		require.Error(t, err, test.name)
	}
}

func TestExtractFileNum(t *testing.T) {
	t.Parallel()

	files, err := ExtractFileNum("1 file")
	require.NoError(t, err)
	require.Equal(t, 1, files)
}

func TestParseIsPublic(t *testing.T) {
	t.Parallel()

	isPublic, err := ParseIsPublic("secret")
	require.NoError(t, err)
	require.False(t, isPublic)
}
