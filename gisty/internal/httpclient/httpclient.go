// Package httpclient builds HTTP clients compatible with GitHub CLI commands.
package httpclient

import (
	"net/http"
	"sort"
	"strings"

	cliapi "github.com/cli/cli/v2/api"
	ghauth "github.com/cli/go-gh/v2/pkg/auth"
	ghconfig "github.com/cli/go-gh/v2/pkg/config"
)

type authTokenGetter struct{}

func (authTokenGetter) ActiveToken(host string) (string, string) {
	return ghauth.TokenForHost(host)
}

func (authTokenGetter) HostForAPIHost(apiHost string) (string, bool) {
	if apiHost == "" {
		return "", false
	}

	cfg, err := ghconfig.Read(nil)
	if err != nil {
		return "", false
	}

	hosts, err := cfg.Keys([]string{"hosts"})
	if err != nil {
		return "", false
	}

	sort.Strings(hosts)

	for _, host := range hosts {
		configured, err := cfg.Get([]string{"hosts", host, "api_host"})
		if err == nil && strings.EqualFold(configured, apiHost) {
			return host, true
		}
	}

	return "", false
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
