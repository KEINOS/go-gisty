// Package httpclient builds HTTP clients compatible with GitHub CLI commands.
package httpclient

import (
	"net/http"

	cliapi "github.com/cli/cli/v2/api"
	ghauth "github.com/cli/go-gh/v2/pkg/auth"
)

type authTokenGetter struct{}

func (authTokenGetter) ActiveToken(host string) (string, string) {
	return ghauth.TokenForHost(host)
}

// New returns a GitHub CLI-compatible HTTP client factory.
func New(appVersion, invokingAgent string) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		return cliapi.NewHTTPClient(cliapi.HTTPClientOptions{
			AppVersion:         appVersion,
			InvokingAgent:      invokingAgent,
			CacheTTL:           0,
			Config:             authTokenGetter{},
			EnableCache:        false,
			Log:                nil,
			LogColorize:        false,
			LogVerboseHTTP:     false,
			SkipDefaultHeaders: false,
			TelemetryDisabler:  nil,
		})
	}
}
