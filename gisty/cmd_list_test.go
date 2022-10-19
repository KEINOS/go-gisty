package gisty

import (
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/gist/list"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGisty_List_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	// mock the list command to force failure
	obj.AltFunctions.List = func(*list.ListOptions) error {
		return errors.New("forced error for listing gists")
	}

	for _, args := range []ListArgs{
		{
			Limit:      1,
			OnlyPublic: true,
			OnlySecret: false, // OnlySecret is prior than OnlyPublic
		},
		{
			Limit:      1,
			OnlyPublic: false,
			OnlySecret: true, // OnlySecret is prior than OnlyPublic
		},
	} {
		// Execute the list command.
		listGistInfo, err := obj.List(args)

		// Assert that the list command failed.
		require.Error(t, err)
		require.Nil(t, listGistInfo, "listGistInfo should be nil on error")
		require.Contains(t, err.Error(), "failed to execute 'gist list' command")
		require.Contains(t, err.Error(), "forced error for listing gists")
	}
}

// ----------------------------------------------------------------------------
//  Private functions
// ----------------------------------------------------------------------------

func Test_parseGistInfo_golden(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		input       string
		gistID      string
		description string
		updatedAt   string
		files       int
		isPublic    bool
	}{
		{
			input:       "d5b9800c636dd78defa4f15894d54d29	Title of gist item2	6 files	secret	2022-04-16T06:08:46Z",
			gistID:      "d5b9800c636dd78defa4f15894d54d29",
			description: "Title of gist item2",
			files:       6,
			isPublic:    false,
			updatedAt:   "2022-04-16 06:08:46 +0000 UTC",
		},
		{
			input:       "\n\ne915aa8c01dd438e3ffd79b05f15a4ff	Title of gist item3	1 file	public	2022-04-18T03:04:38Z",
			gistID:      "e915aa8c01dd438e3ffd79b05f15a4ff",
			description: "Title of gist item3",
			files:       1,
			isPublic:    true,
			updatedAt:   "2022-04-18 03:04:38 +0000 UTC",
		},
		{
			input:       "7101f542be23e5048198e2a27c3cfda8	Title of gist item1	1 file	public	2022-09-18T18:56:10Z\n\n",
			gistID:      "7101f542be23e5048198e2a27c3cfda8",
			description: "Title of gist item1",
			files:       1,
			isPublic:    true,
			updatedAt:   "2022-09-18 18:56:10 +0000 UTC",
		},
	} {
		result, err := parseGistInfo(test.input)
		require.NoError(t, err)
		require.Equal(t, 1, len(result))

		require.Equal(t, test.gistID, result[0].GistID)
		require.Equal(t, test.description, result[0].Description)
		require.Equal(t, test.files, result[0].Files)
		require.Equal(t, test.isPublic, result[0].IsPublic)
		require.Equal(t, test.updatedAt, result[0].UpdatedAt.String())
	}
}

func Test_parseGistInfo_empty_line(t *testing.T) {
	t.Parallel()

	result, err := parseGistInfo("")

	require.NoError(t, err, "parseGistInfo should not return an error on empty line")
	require.Nil(t, result, "result should be nil on empty line")
}

func Test_parseGistInfo_malformed_line(t *testing.T) {
	t.Parallel()

	// updateAt field is missing
	result, err := parseGistInfo("d5b9800c	Title of gist item 2	6 files	secret")

	require.Error(t, err,
		"malformed line should return an error")
	require.Nil(t, result,
		"result should be nil on error")
	require.Contains(t, err.Error(), "failed to parse gist info from:",
		"error message should contain the error reason")
}
