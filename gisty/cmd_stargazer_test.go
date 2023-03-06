package gisty

import (
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/stretchr/testify/require"
)

func TestGisty_Stargazer_msg_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()
	gistID := "5b10b34f87955dfc86d310cd623a61d1"

	obj.AltFunctions.Stargazer = func(*api.ApiOptions) error {
		return NewErr("forced error")
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
		_, err := apiOpt.IO.Out.Write([]byte("unexpected response"))
		require.NoError(t, err)

		return nil
	}

	count, err := obj.Stargazer(gistID)

	require.Error(t, err)
	require.Equal(t, 0, count)
	require.Contains(t, err.Error(), "failed to parse GitHub API response")
	require.Contains(t, err.Error(), "unexpected response")
}
