package httpclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthTokenGetter_ActiveToken(t *testing.T) {
	t.Setenv("GH_TOKEN", "dummy-token")

	token, source := authTokenGetter{}.ActiveToken("github.com")

	require.Equal(t, "dummy-token", token)
	require.Equal(t, "GH_TOKEN", source)
}

func TestNew(t *testing.T) {
	t.Parallel()

	client, err := New("v1.2.3", "go-gisty-test")()

	require.NoError(t, err)
	require.NotNil(t, client)
}
