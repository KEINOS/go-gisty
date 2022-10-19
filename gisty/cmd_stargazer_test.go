package gisty

import (
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGisty_Stargazer_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()
	gistID := "5b10b34f87955dfc86d310cd623a61d1"

	obj.AltFunctions.Stargazer = func(*api.ApiOptions) error {
		return errors.New("forced error")
	}

	count, err := obj.Stargazer(gistID)

	require.Error(t, err)
	require.Equal(t, 0, count)
	require.Contains(t, err.Error(), "failed to execute GitHub API request")
	require.Contains(t, err.Error(), "forced error")
}

func TestGisty_Stargazer_unexpected_response(t *testing.T) {
	t.Parallel()

	obj := NewGisty()
	gistID := "5b10b34f87955dfc86d310cd623a61d1"

	// Success to request but unexpected response.
	obj.AltFunctions.Stargazer = func(apiOpt *api.ApiOptions) error {
		// Set a dummy response.
		apiOpt.IO.Out.Write([]byte("unexpected response"))

		return nil
	}

	count, err := obj.Stargazer(gistID)

	require.Error(t, err)
	require.Equal(t, 0, count)
	require.Contains(t, err.Error(), "failed to parse GitHub API response")
	require.Contains(t, err.Error(), "unexpected response")
}

func Test_sanitizeGistID(t *testing.T) {
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
		actual := sanitizeGistID(test.input)

		require.Equal(t, expect, actual,
			"test #%d: input: %#v", index+1, test.input)
	}
}
